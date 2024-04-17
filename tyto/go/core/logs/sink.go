package logs

type Sink interface {
	// 处理日志
	Handle(record Record)
	// 关闭接收器
	Close()
}
