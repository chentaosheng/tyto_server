package rolling

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
	"tyto/core/fileutil"
	"tyto/core/logging"
	"tyto/core/memutil"
	"tyto/core/osutil"
	"tyto/core/rolling/internal"
	"tyto/core/syncutil"
)

type BufferType = memutil.RefObject[bytes.Buffer]

// 旋转文件writer
type RotateWriter struct {
	options         internal.Options
	fileName        string
	globPattern     string
	writer          internal.Writer
	queue           *syncutil.DoubleQueue[internal.Event]
	nextRotateTime  time.Time
	nextCleanupTime time.Time
	syncTicker      *time.Ticker
	syncDone        chan struct{}
	closed          atomic.Bool
	cleaning        atomic.Bool
	pool            atomic.Pointer[sync.Pool]
	lastError       atomic.Value
}

func NewRotateWriter(opts ...Option) (*RotateWriter, error) {
	o := internal.NewDefaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	// 修正目录格式
	o.OutDir = filepath.Clean(o.OutDir)

	// 检查选项合法性
	if err := o.Validate(); err != nil {
		return nil, err
	}

	// 创建输出目录
	if err := os.MkdirAll(o.OutDir, 0755); err != nil {
		return nil, err
	}

	globPattern := internal.ToGlobPattern(o.NamePattern)

	var (
		nextRotateTime  time.Time
		nextCleanupTime time.Time
	)

	now := time.Now().Local()
	nextRotateTime = internal.GetBaseTime(now, o.RotationInterval)
	nextCleanupTime = internal.GetBaseTime(now, o.CleanupInterval)

	queue := syncutil.NewDoubleQueue[internal.Event](o.WriterQueueSize)

	writer := &RotateWriter{
		options:         *o,
		fileName:        "",
		globPattern:     globPattern,
		writer:          nil,
		queue:           queue,
		nextRotateTime:  nextCleanupTime,
		nextCleanupTime: nextRotateTime,
		syncTicker:      nil,
		syncDone:        nil,
		closed:          atomic.Bool{},
		cleaning:        atomic.Bool{},
		pool:            atomic.Pointer[sync.Pool]{},
		lastError:       atomic.Value{},
	}

	// 启动异步处理协程
	go writer.run()

	return writer, nil
}

func (w *RotateWriter) destroyBuffer(buff *BufferType) {
	w.pool.Load().Put(buff)
}

func (w *RotateWriter) newPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			buff := bytes.Buffer{}
			buff.Grow(128)
			return memutil.NewRefObject(buff, w.destroyBuffer)
		},
	}
}

func (w *RotateWriter) tryInitPool() {
	if w.pool.Load() != nil {
		return
	}

	pool := w.newPool()
	w.pool.CompareAndSwap(nil, pool)
}

// 写入数据，内部会对p进行复制
// 由于是异步写入，返回的err为最近的一次发生的错误
func (w *RotateWriter) Write(p []byte) (n int, err error) {
	if w.closed.Load() {
		return 0, os.ErrClosed
	}

	w.tryInitPool()

	buff := w.pool.Load().Get().(*BufferType)
	buff.IncRef()
	buff.Object().Reset()
	buff.Object().Write(p)

	return w.writeBuffer(buff)
}

// 写入数据，内部会对s进行复制
// 由于是异步写入，返回的err为最近的一次发生的错误
func (w *RotateWriter) WriteString(s string) (n int, err error) {
	if w.closed.Load() {
		return 0, os.ErrClosed
	}

	w.tryInitPool()

	buff := w.pool.Load().Get().(*BufferType)
	buff.IncRef()
	buff.Object().Reset()
	buff.Object().WriteString(s)

	return w.writeBuffer(buff)
}

func (w *RotateWriter) writeBuffer(buff *BufferType) (n int, err error) {
	event := internal.Event{
		Type:   internal.EVENT_TYPE_BUFFER,
		Buffer: buff,
	}

	n = buff.Object().Len()

	v := w.lastError.Load()
	if v != nil {
		err = v.(error)
		w.lastError.Store(nil)
	}

	// push后，buff的所有权交个另一个go routine
	// 不要再使用buff
	w.queue.Push(event)

	return
}

func (w *RotateWriter) WriteBuffer(buff *BufferType) (n int, err error) {
	if w.closed.Load() {
		return 0, os.ErrClosed
	}

	buff.IncRef()

	return w.writeBuffer(buff)
}

func (w *RotateWriter) sync() error {
	if w.closed.Load() {
		return nil
	}

	// writer为空，无需同步
	if w.writer == nil {
		return nil
	}

	// 通常是操作者手动删除了文件，导致writer失效
	// 因此，关闭writer，等待下次写入时重新创建
	// 这种情况下，会丢失一定时间内的数据
	// 为了性能，不在每次写入时检查writer是否有效
	if !w.writer.IsValid() {
		_ = w.writer.Close()
		w.writer = nil
		return nil
	}

	return w.writer.Sync()
}

// 同步阻塞
func (w *RotateWriter) Sync() error {
	if w.closed.Load() {
		return nil
	}

	c := make(chan error)
	event := internal.Event{
		Type: internal.EVENT_TYPE_SYNC,
		Chan: c,
	}

	w.queue.Push(event)

	// wait
	err := <-c
	close(c)

	return err
}

func (w *RotateWriter) cleanupQueue() {
	// 确保所有event都被处理
	for !w.queue.Empty() {
		e := w.queue.Pop()

		switch e.Type {
		case internal.EVENT_TYPE_SYNC:
			if e.Chan == nil {
				continue
			}
			e.Chan <- nil

		case internal.EVENT_TYPE_CLOSE:
			e.Chan <- nil

		case internal.EVENT_TYPE_BUFFER:
			err := w.handleBuffer(e.Buffer)
			if err != nil {
				w.Logger().Error("write buffer failed when cleanup:", err.Error())
			}
		}
	}
}

func (w *RotateWriter) close() error {
	if w.closed.Load() {
		return nil
	}
	if !w.closed.CompareAndSwap(false, true) {
		return nil
	}

	if w.writer == nil && w.queue.Empty() {
		return nil
	}

	// 确保所有event都被处理
	w.cleanupQueue()

	if w.writer == nil {
		return nil
	}

	// 关闭writer
	err := w.writer.Close()
	w.writer = nil
	return err
}

// 同步阻塞
func (w *RotateWriter) Close() error {
	if w.closed.Load() {
		return nil
	}

	c := make(chan error)
	event := internal.Event{
		Type: internal.EVENT_TYPE_CLOSE,
		Chan: c,
	}

	w.queue.Push(event)

	// wait
	err := <-c
	close(c)

	return err
}

func (w *RotateWriter) handleBuffer(buff *BufferType) error {
	var err error

	if err = w.rotate(); err != nil {
		// 通常是权限问题、磁盘空间问题导致文件创建失败
		return err
	}

	_, err = w.writer.Write(buff.Object().Bytes())

	// 释放资源
	buff.DecRef()

	return err
}

func (w *RotateWriter) stopTicker() {
	if w.syncTicker == nil {
		return
	}

	w.syncTicker.Stop()
	w.syncDone <- struct{}{}
	close(w.syncDone)
}

func (w *RotateWriter) handleSyncTicker() {
	for {
		select {
		case <-w.syncTicker.C:
			if w.closed.Load() {
				continue
			}

			// 不需要chan获取结果
			event := internal.Event{
				Type: internal.EVENT_TYPE_SYNC,
				Chan: nil,
			}

			w.queue.Push(event)

		case <-w.syncDone:
			return
		}
	}
}

func (w *RotateWriter) run() {
	defer w.stopTicker()

	for {
		e := w.queue.Pop()

		switch e.Type {
		case internal.EVENT_TYPE_SYNC:
			err := w.sync()
			if e.Chan != nil {
				e.Chan <- err

			} else if err != nil {
				w.lastError.Store(err)
			}

		case internal.EVENT_TYPE_CLOSE:
			err := w.close()
			e.Chan <- err
			// 退出
			return

		case internal.EVENT_TYPE_BUFFER:
			err := w.handleBuffer(e.Buffer)
			if err != nil {
				w.lastError.Store(err)
			}
		}
	}
}

func (w *RotateWriter) rotate() error {
	now := time.Now().Local()
	if now.Before(w.nextRotateTime) && w.writer != nil {
		// 不需要旋转
		return nil
	}

	baseTime := internal.GetBaseTime(now, w.options.RotationInterval)
	fileName := internal.GenerateFileName(w.options.NamePattern, baseTime)

	if w.fileName == fileName && w.writer != nil {
		// 当前文件名与新文件名相同，不需要旋转
		w.nextRotateTime = baseTime.Add(w.options.RotationInterval)
		return nil
	}

	// 创建新文件
	fullPath := filepath.Join(w.options.OutDir, fileName)
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// 更新writer
	if w.writer != nil {
		w.writer.Reset(file)
	} else {
		w.writer = internal.NewWriter(file, w.options.BufferSize)
	}

	// 开启定时刷新缓存功能
	if w.writer.IsBuffered() && w.syncTicker == nil {
		w.syncTicker = time.NewTicker(w.options.FlushInterval)
		w.syncDone = make(chan struct{})
		go w.handleSyncTicker()
	}

	// 更新旋转信息
	w.fileName = fileName
	w.nextRotateTime = baseTime.Add(w.options.RotationInterval)

	// 创建软链接
	if (osutil.IsLinux() || osutil.IsMacOsx()) && len(w.options.LinkName) != 0 {
		linkPath := filepath.Join(w.options.OutDir, w.options.LinkName)
		tempPath := linkPath + ".link"

		err = os.Symlink(w.fileName, tempPath)
		// 1. 失败也没那么重要
		// 2. 一般不会出错
		if err != nil {
			w.Logger().Error("symlink", tempPath, "->", w.fileName, "failed:", err.Error())

		} else {
			err = os.Rename(tempPath, linkPath)
			if err != nil {
				w.Logger().Error("rename", tempPath, "->", linkPath, "failed:", err.Error())
			}
		}
	}

	// 清理过时文件，每次旋转都会尝试清理一次
	if w.options.MaxAge > 0 && now.After(w.nextCleanupTime) {
		baseTime = internal.GetBaseTime(now, w.options.CleanupInterval)
		w.nextCleanupTime = baseTime.Add(w.options.CleanupInterval)
		w.cleanup()
	}

	return nil
}

// 清理过时的文件
func (w *RotateWriter) cleanup() {
	// 文件清理过于耗时，清理间隔过短，可能会导致同时启动多个清理协程，因而进行限制
	if w.cleaning.Load() {
		return
	}
	if !w.cleaning.CompareAndSwap(false, true) {
		return
	}

	fileList, err := fileutil.ListFile(w.options.OutDir, w.globPattern)
	if err != nil {
		w.cleaning.Store(false)
		w.Logger().Error("list file failed:", err.Error())
		return
	}

	currentFile := filepath.Join(w.options.OutDir, w.fileName)
	expiredTime := time.Now().Local().Add(-w.options.MaxAge)

	go func() {
		defer w.cleaning.Store(false)

		for _, file := range fileList {
			// 跳过当前文件
			if file == currentFile {
				continue
			}

			fileInfo, err := os.Lstat(file)
			if err != nil {
				continue
			}

			if !fileInfo.Mode().IsRegular() {
				// 不是正常的文件
				continue
			}

			if fileInfo.ModTime().After(expiredTime) {
				// 未过期
				continue
			}

			err = os.Remove(file)
			if err != nil {
				w.Logger().Error("remove", file, "failed:", err.Error())
			}
		}
	}()
}

func (w *RotateWriter) Logger() logging.Logger {
	return w.options.Logger
}
