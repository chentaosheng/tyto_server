package strutil

import (
	"fmt"
	"strconv"
	"strings"
	"tyto/core/sqlutil"
	"tyto/core/tyto"
)

type StringBuilder struct {
	strings.Builder
	cache [32]byte
}

func NewStringBuilder(size int) *StringBuilder {
	sb := &StringBuilder{}
	sb.Grow(size)
	return sb
}

func (sb *StringBuilder) WriteInt8(i int8) {
	sb.WriteInt64(int64(i))
}

func (sb *StringBuilder) WriteInt16(i int16) {
	sb.WriteInt64(int64(i))
}

func (sb *StringBuilder) WriteInt32(i int32) {
	sb.WriteInt64(int64(i))
}

func (sb *StringBuilder) WriteInt64(i int64) {
	cache := sb.cache[:0]
	strconv.AppendInt(cache, i, 10)
	sb.Write(cache)
}

func (sb *StringBuilder) WriteUint8(n uint8) {
	sb.WriteUint64(uint64(n))
}

func (sb *StringBuilder) WriteUint16(n uint16) {
	sb.WriteUint64(uint64(n))
}

func (sb *StringBuilder) WriteUint32(n uint32) {
	sb.WriteUint64(uint64(n))
}

func (sb *StringBuilder) WriteUint64(n uint64) {
	cache := sb.cache[:0]
	strconv.AppendUint(cache, n, 10)
	sb.Write(cache)
}

func (sb *StringBuilder) writeFloat(f float64, prec int, bitSize int) {
	cache := sb.cache[:0]
	strconv.AppendFloat(cache, f, 'f', prec, bitSize)
	sb.Write(cache)
}

// prec 表示小数点后的位数
func (sb *StringBuilder) WriteFloat32(f float32, prec int32) {
	sb.writeFloat(float64(f), int(prec), 32)
}

// prec 表示小数点后的位数
func (sb *StringBuilder) WriteFloat64(f float64, prec int32) {
	sb.writeFloat(f, int(prec), 64)
}

func (sb *StringBuilder) WriteBool(b bool) {
	if b {
		sb.WriteString("true")
	} else {
		sb.WriteString("false")
	}
}

func (sb *StringBuilder) EscapeWrite(bs []byte) {
	sb.Write(sqlutil.EscapeBytes(bs))
}

func (sb *StringBuilder) EscapeWriteString(s string) {
	sb.WriteString(sqlutil.EscapeString(s))
}

func (sb *StringBuilder) Format(ctx tyto.Context, format string, a ...interface{}) bool {
	if _, err := fmt.Fprintf(sb, format, a...); err != nil {
		ctx.Logger().Error("format failed:", err)
		return false
	}
	return true
}
