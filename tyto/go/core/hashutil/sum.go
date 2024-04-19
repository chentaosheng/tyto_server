package hashutil

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"tyto/core/tyto"
)

func hashSum(ctx tyto.Context, h hash.Hash, format string, a ...interface{}) string {
	buff := bytes.NewBuffer(make([]byte, 0, 256))
	if _, err := fmt.Fprintf(buff, format, a...); err != nil {
		ctx.Logger().Error("format failed:", err)
		return ""
	}

	h.Write(buff.Bytes())
	sum := h.Sum(nil)

	return hex.EncodeToString(sum)
}

func hashSumToBase64(ctx tyto.Context, h hash.Hash, format string, a ...interface{}) string {
	buff := bytes.NewBuffer(make([]byte, 0, 256))
	if _, err := fmt.Fprintf(buff, format, a...); err != nil {
		ctx.Logger().Error("format failed:", err)
		return ""
	}

	h.Write(buff.Bytes())
	sum := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(sum)
}
