#include <chrono>
#include <thread>
#include <cmath>
#include <random>
#include "cpu_relax.hpp"
#include "spin_lock.hpp"

namespace tyto
{
	void SpinLock::Lock() noexcept
	{
		// 自旋方式等待的次数
		constexpr std::size_t kPauseCountBeforeSleep = 32;
		constexpr std::size_t kSleepCountBeforeYield = 64;
		// 计算退避时，用于计算退避次数的偏移值
		constexpr std::size_t kMaxBackOffShift = 16;

		static thread_local std::minstd_rand generator{ std::random_device{}() };
		// 获取锁时，与其他线程冲突的次数
		std::size_t conflict_count = 0;

		for (;;)
		{
			std::size_t retries = 0;
			while (flag_.test(std::memory_order_relaxed))
			{
				// 按顺序用不同的方式自旋等待
				if (retries < kPauseCountBeforeSleep)
				{
					++retries;
					CpuRelax();
				}
				else if (retries < kSleepCountBeforeYield)
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

			if (flag_.test_and_set(std::memory_order_acquire))
			{
				// 获取失败，说明线程间竞争比较激烈，进行随机退避
				std::size_t shift = (std::min)(conflict_count, kMaxBackOffShift);
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
}
