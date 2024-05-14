#pragma once

#include <boost/circular_buffer.hpp>
#include "typedef.hpp"

namespace tyto
{
	// ID池
	class IdPool
	{
	public:
		// 无效ID值
		static constexpr int32 kInvalidId = 0;

	public:
		// max_id: ID的最大值
		// min_free_count: 最小空闲ID数量，空闲列表中的ID超过此数量则直接重用旧ID
		IdPool(int32 max_id, int32 min_free_count);

		// 获取一个ID，范围[1, max_id)
		int32 Get();
		// 回收ID
		void Put(int32 id);

	private:
		bool Exhausted();
		int32 Reuse();

	private:
		int32 next_id_;
		int32 max_id_;
		int32 min_free_count_;
		boost::circular_buffer_space_optimized<int32> free_list_;
	};
}