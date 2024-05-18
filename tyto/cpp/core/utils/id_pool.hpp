#pragma once

#include <optional>
#include <type_traits>
#include <boost/circular_buffer.hpp>
#include "typedef.hpp"

namespace tyto
{
	template<typename T, typename Enable = void>
	class IdPool
	{
	};

	// ID池
	// 注意，如果ID范围太大，将可能占用大量内存
	template<typename T>
	class IdPool<T, typename std::enable_if_t<std::is_integral_v<T>>>
	{
		// 无效ID值
		static constexpr T kInvalidId = 0;

	public:
		// max_id: ID的最大值
		// min_free: 最小空闲ID数量，空闲列表中的ID超过此数量则直接重用旧ID
		// 若min_free为0，则用完所有id开始使用回收的ID，但这样会占用更多内存
		explicit IdPool(T max_id, T min_free);

		// 获取一个ID，范围[1, max_id)
		// 若没有足够
		std::optional<T> Get();
		// 回收ID
		void Put(T id);

	private:
		// ID是否已经耗尽
		bool Exhausted();
		// 复用ID
		T Reuse();

	private:
		T next_id_;
		T max_id_;
		T min_free_;
		boost::circular_buffer_space_optimized<T> free_list_;
	};

	template<typename T>
	IdPool<T, typename std::enable_if_t<std::is_integral_v<T>>>::IdPool(T max_id, T min_free)
		: next_id_(1)
		, max_id_(max_id)
		, min_free_(min_free)
		, free_list_(max_id)
	{
	}

	template<typename T>
	std::optional<T> IdPool<T, typename std::enable_if_t<std::is_integral_v<T>>>::Get()
	{
		// 没有可分配的ID
		if (Exhausted())
			return std::nullopt;

		// 没有可分配的ID，但有空闲的ID
		if (next_id_ >= max_id_)
			return Reuse();

		// 超过最小空闲id数量，直接重用旧id
		if (min_free_ > 0 && free_list_.size() > min_free_)
			return Reuse();

		return next_id_++;
	}

	template<typename T>
	void IdPool<T, typename std::enable_if_t<std::is_integral_v<T>>>::Put(T id)
	{
		if (id <= kInvalidId || id >= max_id_)
			return;

		free_list_.push_back(id);
	}

	template<typename T>
	bool IdPool<T, typename std::enable_if_t<std::is_integral_v<T>>>::Exhausted()
	{
		return next_id_ >= max_id_ && free_list_.empty();
	}

	template<typename T>
	T IdPool<T, typename std::enable_if_t<std::is_integral_v<T>>>::Reuse()
	{
		T id = free_list_.front();
		free_list_.pop_front();
		return id;
	}
}
