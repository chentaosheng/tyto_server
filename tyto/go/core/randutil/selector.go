package randutil

import "tyto/core/tyto"

type Selector[T any] interface {
	// 添加选择项
	Add(ctx tyto.Context, value T, weight int32)
	// 构建选择列表
	BuildList()
	// 随机选择一个元素，可重复选择
	Select(ctx tyto.Context) (T, bool)
	// 随机选择一个元素，不可重复选择，首次调用前需要调用Reset()以前的选择记录
	UniqueSelect(ctx tyto.Context) (T, bool)
	// 清空选择记录
	Reset()
}
