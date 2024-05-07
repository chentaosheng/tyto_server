#pragma once

#include <atomic>
#include <boost/noncopyable.hpp>

namespace tyto
{
	//
	// 自旋锁
	class SpinLock : public boost::noncopyable
	{
	public:
		SpinLock() = default;

		// 加锁
		void Lock() noexcept;

		// 尝试加锁，非阻塞，加锁成功返回true，否则返回false
		bool TryLock() noexcept
		{
			return !flag_.test_and_set(std::memory_order_acquire);
		}

		// 解锁
		void Unlock() noexcept
		{
			flag_.clear(std::memory_order_release);
		}

	private:
		std::atomic_flag flag_;
	};
}
