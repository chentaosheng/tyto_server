package randutil

// 判断是否在几率范围内，用于命中、暴击判断等
// odds 发生的概率
// max  概率的最大值
func InOdds(odds int32, max int32) bool {
	return Int32R(0, max) < odds
}
