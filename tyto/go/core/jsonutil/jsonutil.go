package jsonutil

import (
	"encoding/json"
	"tyto/core/byteutil"
	"tyto/core/sqlutil"
	"tyto/core/tyto"
)

// struct -> []byte
func Marshal(ctx tyto.Context, v interface{}) ([]byte, bool) {
	bs, err := json.Marshal(v)
	if err != nil {
		ctx.Logger().Error("marshal failed:", err)
		return []byte("{}"), false
	}

	return bs, true
}

// []byte -> struct
func Unmarshal(ctx tyto.Context, bs []byte, v interface{}) bool {
	if err := json.Unmarshal(bs, v); err != nil {
		ctx.Logger().Error("unmarshal failed:", err)
		return false
	}

	return true
}

// struct -> string
func MarshalToString(ctx tyto.Context, v interface{}) (string, bool) {
	bs, ok := Marshal(ctx, v)
	return byteutil.BytesToStr(bs), ok
}

// string -> struct
func UnmarshalFromString(ctx tyto.Context, s string, v interface{}) bool {
	return Unmarshal(ctx, byteutil.StrToBytes(s), v)
}

// struct -> []byte
// 有转义，结果可直接作为sql语句的value值
func EscapeMarshal(ctx tyto.Context, v interface{}) ([]byte, bool) {
	bs, err := json.Marshal(v)
	if err != nil {
		ctx.Logger().Error("marshal failed:", err)
		return []byte("{}"), false
	}

	return sqlutil.EscapeBytes(bs), true
}

// struct -> string
// 有转义，结果可直接作为sql语句的value值
func EscapeMarshalToString(ctx tyto.Context, v interface{}) (string, bool) {
	bs, ok := EscapeMarshal(ctx, v)
	return byteutil.BytesToStr(bs), ok
}

// 将数据库查询结果转化为struct
func UnmarshalFromDB(ctx tyto.Context, res interface{}, v interface{}) bool {
	// pointer to bytes
	pbs, ok := res.(*[]byte)
	if !ok {
		ctx.Logger().Error("invalid type")
		return false
	}

	if pbs == nil || len(*pbs) == 0 {
		empty := []byte("{}")
		pbs = &empty
	}

	return Unmarshal(ctx, *pbs, v)
}
