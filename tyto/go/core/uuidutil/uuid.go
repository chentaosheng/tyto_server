package uuidutil

import (
	"encoding/hex"
	"github.com/google/uuid"
	"tyto/core/tyto"
)

func init() {
	// 提高生成uuid的速度
	uuid.EnableRandPool()
}

// 新建uuid字符串，32字节，格式：018e414e50f27b3fbe6492fcb1c21a00
func NewUuidString(ctx tyto.Context) (string, bool) {
	id, err := uuid.NewV7()
	if err != nil {
		ctx.Logger().Error("new uuid failed:", err.Error())
		return "", false
	}

	b, err := id.MarshalBinary()
	if err != nil {
		ctx.Logger().Error("marshal uuid failed:", err.Error())
		return "", false
	}

	return hex.EncodeToString(b), true
}
