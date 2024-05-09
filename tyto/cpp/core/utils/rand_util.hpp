#pragma once

#include "typedef.hpp"

namespace tyto
{
	class RandUtil
	{
	public:
		// 随机生成[min, max)之间的整数
		static int32 Int32R(int32 min, int32 max);
		static int64 Int64R(int64 min, int64 max);

		// 随机生成[min, max)之间的整数
		static uint32 Uint32R(uint32 min, uint32 max);
		static uint64 Uint64R(uint64 min, uint64 max);

		// 随机生成[0, max)之间的整数
		static int32 Int32N(int32 max);
		static int64 Int64N(int64 max);

		// 随机生成[0, max)之间的整数
		static uint32 Uint32N(uint32 max);
		static uint64 Uint64N(uint64 max);

	public:
		RandUtil() = delete;
	};
}
