#include <chrono>
#include "time_util.hpp"

namespace tyto
{
	int64 TimeUtil::NowMilli()
	{
		auto now = std::chrono::steady_clock::now();
		auto duration = now.time_since_epoch();
		return std::chrono::duration_cast<std::chrono::milliseconds>(duration).count();
	}

	int64 TimeUtil::NowSecond()
	{
		auto now = std::chrono::steady_clock::now();
		auto duration = now.time_since_epoch();
		return std::chrono::duration_cast<std::chrono::seconds>(duration).count();
	}
}
