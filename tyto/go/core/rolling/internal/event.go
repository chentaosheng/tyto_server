package internal

type EventType int32

const (
	EVENT_TYPE_SYNC  EventType = 1
	EVENT_TYPE_CLOSE EventType = 2
)

type Event struct {
	Type EventType  // 事件类型
	C    chan error // 返回处理结果用的channel，可能为nil
}
