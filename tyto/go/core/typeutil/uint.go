package typeutil

// 用high和low组成一个新的uint，high放高位，low放低位
func MakeUint16(high uint8, low uint8) uint16 {
	data := uint16(high) << 8
	data |= uint16(low)
	return data
}

func MakeUint32(high uint16, low uint16) uint32 {
	data := uint32(high) << 16
	data |= uint32(low)
	return data
}

func MakeUint64(high uint32, low uint32) uint64 {
	data := uint64(high) << 32
	data |= uint64(low)
	return data
}

// 获取高位部分
func HighPartUint16(data uint16) uint8 {
	data >>= 8
	return uint8(data)
}

func HighPartUint32(data uint32) uint16 {
	data >>= 16
	return uint16(data)
}

func HighPartUint64(data uint64) uint32 {
	data >>= 32
	return uint32(data)
}

// 获取低位部分
func LowPartUint16(data uint16) uint8 {
	data &= 0xff
	return uint8(data)
}

func LowPartUint32(data uint32) uint16 {
	data &= 0xffff
	return uint16(data)
}

func LowPartUint64(data uint64) uint32 {
	data &= 0xffffffff
	return uint32(data)
}
