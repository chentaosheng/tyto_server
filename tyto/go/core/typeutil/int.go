package typeutil

// 用high和low组成一个新的int，high放高位，low放低位
func MakeInt16(high int8, low int8) int16 {
	// int8(-1) 和 int16(-1) 的二进制数据是不同的
	// 因此，先转换为同长度的无符号类型
	data := uint16(uint8(high)) << 8
	data |= uint16(uint8(low))
	return int16(data)
}

func MakeInt32(high int16, low int16) int32 {
	data := uint32(uint16(high)) << 16
	data |= uint32(uint16(low))
	return int32(data)
}

func MakeInt64(high int32, low int32) int64 {
	data := uint64(uint32(high)) << 32
	data |= uint64(uint32(low))
	return int64(data)
}

// 获取高位部分
func HighPartInt16(i int16) int8 {
	data := uint16(i) >> 8
	return int8(data)
}

func HighPartInt32(i int32) int16 {
	data := uint32(i) >> 16
	return int16(data)
}

func HighPartInt64(i int64) int32 {
	data := uint64(i) >> 32
	return int32(data)
}

// 获取低位部分
func LowPartInt16(i int16) int8 {
	data := uint16(i) & 0xff
	return int8(data)
}

func LowPartInt32(i int32) int16 {
	data := uint32(i) & 0xffff
	return int16(data)
}

func LowPartInt64(i int64) int32 {
	data := uint64(i) & 0xffffffff
	return int32(data)
}
