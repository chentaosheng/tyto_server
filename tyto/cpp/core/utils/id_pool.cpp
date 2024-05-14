#include "id_pool.hpp"

namespace tyto
{
	IdPool::IdPool(int32 max_id, int32 min_free_count)
		: next_id_(1)
		, max_id_(max_id)
		, min_free_count_(min_free_count)
		, free_list_(max_id)
	{

	}

	int32 IdPool::Get()
	{
		// 没有可分配的ID
		if (Exhausted())
			return kInvalidId;

		// 没有可分配的ID，但有空闲的ID
		if (next_id_ >= max_id_)
			return Reuse();

		// 超过最小空闲id数量，直接重用旧id
		if (min_free_count_ > 0 && free_list_.size() > min_free_count_)
			return Reuse();

		return next_id_++;
	}

	void IdPool::Put(int32 id)
	{
		if (id <= kInvalidId || id >= max_id_)
			return;

		free_list_.push_back(id);
	}

	bool IdPool::Exhausted()
	{
		return next_id_ >= max_id_ && free_list_.empty();
	}

	int32 IdPool::Reuse()
	{
		int32 id = free_list_.front();
		free_list_.pop_front();
		return id;
	}
}