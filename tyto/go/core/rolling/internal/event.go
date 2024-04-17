package internal

import (
	"bytes"
	"tyto/core/memutil"
)

type EventType int32

const (
	EVENT_TYPE_SYNC   EventType = 1
	EVENT_TYPE_CLOSE  EventType = 2
	EVENT_TYPE_BUFFER EventType = 3
)

type Event struct {
	Type   EventType                        // 事件类型
	Chan   chan error                       // 返回处理结果用的channel，可能为nil
	Buffer *memutil.RefObject[bytes.Buffer] // 需要写入的数据
}
