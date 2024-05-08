#pragma once

#include "typedef.hpp"

namespace tyto
{
	// 时间工具类
	class TimeUtil
	{
	public:
		// 获取当前时间戳，单位为毫秒
		static int64 NowMilli();

		// 获取当前时间戳，单位为秒
		static int64 NowSecond();

	public:
		TimeUtil() = delete;
	};
}
