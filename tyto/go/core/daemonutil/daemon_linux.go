package daemonutil

import (
	"fmt"
	"github.com/sevlyar/go-daemon"
	"golang.org/x/exp/slices"
	"os"
	"os/exec"
	"path/filepath"
	"tyto/core/logging"
)

// 变为守护进程
func Start(logger logging.Logger, name string, code string) (*daemon.Context, bool) {
	// 获取exe的路径
	// 该路径不一定是一个真实存在的文件，可能是一个软链接，也可能exe文件已经被删除
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		logger.Error("look exe path failed:", err)
		return nil, false
	}

	// exe的绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		logger.Error("get exe absolute path failed:", err)
		return nil, false
	}

	// pid文件
	pidDir := "./run"
	err = os.MkdirAll(pidDir, 0755)
	if err != nil {
		logger.Error("create pid dir failed:", err)
		return nil, false
	}
	pidPath := fmt.Sprintf("%s/%s-%s.pid", pidDir, name, code)

	// stdout/stderr输出到文件
	outDir := "./out"
	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		logger.Error("create out dir failed:", err)
		return nil, false
	}
	outPath := fmt.Sprintf("%s/%s-%s.out", outDir, name, code)

	// 修改路径，方便linux用ps等命令查看
	args := slices.Clone(os.Args)
	args[0] = absPath

	dctx := &daemon.Context{
		PidFileName: pidPath,
		PidFilePerm: 0644,
		LogFileName: outPath,
		LogFilePerm: 0644,
		WorkDir:     "./",
		Umask:       022,
		Args:        args,
	}

	// 模仿fork()的逻辑
	child, err := dctx.Reborn()
	if err != nil {
		logger.Error("daemon reborn failed:", err)
		return nil, false
	}

	if child != nil {
		// 父进程退出
		os.Exit(0)
	}

	// 子进程继续执行
	return dctx, true
}

func Stop(logger logging.Logger, dctx *daemon.Context) bool {
	logger.Error("stop daemon failed, system unsupported")
	return false
}
