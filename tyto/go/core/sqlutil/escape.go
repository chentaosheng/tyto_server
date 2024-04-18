package sqlutil

import (
	"bytes"
	"strings"
)

// 对特殊字符进行转义，以便安全地插入到SQL语句中
// see `mysql-server/mysys/charset.cc`
func EscapeString(str string) string {
	sb := strings.Builder{}
	sb.Grow(len(str) * 2)

	for i := 0; i < len(str); i++ {
		c := str[i]
		switch c {
		case 0:
			sb.WriteByte('\\')
			sb.WriteByte('0')
		case '\n':
			sb.WriteByte('\\')
			sb.WriteByte('n')
		case '\r':
			sb.WriteByte('\\')
			sb.WriteByte('r')
		case '\\':
			sb.WriteByte('\\')
			sb.WriteByte('\\')
		case '\'':
			sb.WriteByte('\\')
			sb.WriteByte('\'')
		case '"':
			sb.WriteByte('\\')
			sb.WriteByte('"')
		case '\032':
			sb.WriteByte('\\')
			sb.WriteByte('Z')
		default:
			sb.WriteByte(c)
		}
	}

	return sb.String()
}

// 对特殊字符进行转义，以便安全地插入到SQL语句中
func EscapeBytes(bs []byte) []byte {
	buff := bytes.NewBuffer(make([]byte, 0, len(bs)*2))

	for i := 0; i < len(bs); i++ {
		c := bs[i]
		switch c {
		case 0:
			buff.WriteByte('\\')
			buff.WriteByte('0')
		case '\n':
			buff.WriteByte('\\')
			buff.WriteByte('n')
		case '\r':
			buff.WriteByte('\\')
			buff.WriteByte('r')
		case '\\':
			buff.WriteByte('\\')
			buff.WriteByte('\\')
		case '\'':
			buff.WriteByte('\\')
			buff.WriteByte('\'')
		case '"':
			buff.WriteByte('\\')
			buff.WriteByte('"')
		case '\032':
			buff.WriteByte('\\')
			buff.WriteByte('Z')
		default:
			buff.WriteByte(c)
		}
	}

	return buff.Bytes()
}
