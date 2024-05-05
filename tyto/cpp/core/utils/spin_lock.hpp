#pragma once

#include <atomic>
#include <chrono>
#include <thread>
#include <cmath>
#include <random>
#include <boost/noncopyable.hpp>
#include "cpu_relax.hpp"

namespace tyto
{
	//
	// 自旋锁
	class SpinLock : public boost::noncopyable
	{
	public:
		SpinLock() = default;

		// 加锁
		void Lock() noexcept
		{
			static_assert(flag_.is_always_lock_free);

			// 自旋方式等待的次数
			constexpr int kPAUSE_COUNT_BEFORE_SLEEP = 32;
			constexpr int kSLEEP_COUNT_BEFORE_YIELD = 64;
			// 计算退避时，用于计算退避次数的偏移值
			constexpr std::size_t kMAX_BACKOFF_SHIFT = 16;

			static thread_local std::minstd_rand generator{ std::random_device{}() };
			// 获取锁时，与其他线程冲突的次数
			std::size_t conflict_count = 0;

			for (;;)
			{
				std::size_t retries = 0;
				while (flag_.load(std::memory_order_relaxed))
				{
					// 按顺序用不同的方式自旋等待
					if (retries < kPAUSE_COUNT_BEFORE_SLEEP)
					{
						++retries;
						CpuRelax();
					}
					else if (retries < kSLEEP_COUNT_BEFORE_YIELD)
					{
						++retries;
						static constexpr std::chrono::microseconds sleep_time{ 0 };
						std::this_thread::sleep_for(sleep_time);
					}
					else
					{
						std::this_thread::yield();
					}
				}

				if (!flag_.exchange(true, std::memory_order_acquire))
				{
					// 获取失败，说明线程间竞争比较激烈，进行随机退避
					std::size_t shift = (std::min)(conflict_count, kMAX_BACKOFF_SHIFT);
					std::uniform_int_distribution<std::size_t> dist{
						0,
						static_cast<std::size_t>(1) << shift
					};

					const std::size_t relax_count = dist(generator);
					++conflict_count;

					for (std::size_t i = 0; i < relax_count; ++i)
					{
						CpuRelax();
					}

					continue;
				}
				else
				{
					// 加锁成功
					break;
				}
			}
		}

		// 尝试加锁，非阻塞，加锁成功返回true，否则返回false
		bool TryLock() noexcept
		{
			return !flag_.exchange(true, std::memory_order_acquire);
		}

		// 解锁
		void Unlock() noexcept
		{
			flag_.store(false, std::memory_order_release);
		}

	private:
		std::atomic_bool flag_;
	};
}
