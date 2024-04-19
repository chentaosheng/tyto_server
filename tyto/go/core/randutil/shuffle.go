package randutil

// 打乱切片的元素顺序
func Shuffle[T any](slice []T) {
	for i := len(slice) - 1; i > 0; i-- {
		j := Int32N(int32(i + 1))
		slice[i], slice[j] = slice[j], slice[i]
	}
}
