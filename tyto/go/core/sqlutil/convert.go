package sqlutil

import (
	"strconv"
	"tyto/core/byteutil"
	"tyto/core/tyto"
)

// 查询数据库时，可为null的字段，会返回[]byte类型数据，
// 因此，需要通过下面的方法转化为int/uint/float类型
func toInt(ctx tyto.Context, bs []byte, bitSize int) (int64, bool) {
	if len(bs) == 0 {
		return 0, true
	}

	str := byteutil.BytesToStr(bs)
	v, err := strconv.ParseInt(str, 10, bitSize)
	if err != nil {
		ctx.Logger().Error("parse int failed:", err)
		return 0, false
	}

	return v, true
}

func ToInt8(ctx tyto.Context, bs []byte) (int8, bool) {
	v, ok := toInt(ctx, bs, 8)
	return int8(v), ok
}

func ToInt16(ctx tyto.Context, bs []byte) (int16, bool) {
	v, ok := toInt(ctx, bs, 16)
	return int16(v), ok
}

func ToInt32(ctx tyto.Context, bs []byte) (int32, bool) {
	v, ok := toInt(ctx, bs, 32)
	return int32(v), ok
}

func ToInt64(ctx tyto.Context, bs []byte) (int64, bool) {
	return toInt(ctx, bs, 64)
}

func toUint(ctx tyto.Context, bs []byte, bitSize int) (uint64, bool) {
	if len(bs) == 0 {
		return 0, true
	}

	str := byteutil.BytesToStr(bs)
	v, err := strconv.ParseUint(str, 10, bitSize)
	if err != nil {
		ctx.Logger().Error("parse uint failed:", err)
		return 0, false
	}

	return v, true
}

func ToUint8(ctx tyto.Context, bs []byte) (uint8, bool) {
	v, ok := toUint(ctx, bs, 8)
	return uint8(v), ok
}

func ToUint16(ctx tyto.Context, bs []byte) (uint16, bool) {
	v, ok := toUint(ctx, bs, 16)
	return uint16(v), ok
}

func ToUint32(ctx tyto.Context, bs []byte) (uint32, bool) {
	v, ok := toUint(ctx, bs, 32)
	return uint32(v), ok
}

func ToUint64(ctx tyto.Context, bs []byte) (uint64, bool) {
	return toUint(ctx, bs, 64)
}

func toFloat(ctx tyto.Context, bs []byte, bitSize int) (float64, bool) {
	if len(bs) == 0 {
		return 0, true
	}

	str := byteutil.BytesToStr(bs)
	v, err := strconv.ParseFloat(str, bitSize)
	if err != nil {
		ctx.Logger().Error("parse float failed:", err)
		return 0, false
	}

	return v, true
}

func ToFloat32(ctx tyto.Context, bs []byte) (float32, bool) {
	v, ok := toFloat(ctx, bs, 32)
	return float32(v), ok
}

func ToFloat64(ctx tyto.Context, bs []byte) (float64, bool) {
	return toFloat(ctx, bs, 64)
}
