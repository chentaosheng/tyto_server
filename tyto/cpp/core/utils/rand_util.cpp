#include <random>
#include "rand_util.hpp"

namespace tyto
{
	static std::mt19937& GetRandomEngine32()
	{
		static thread_local std::mt19937 gen(std::random_device{}());
		return gen;
	}

	static std::mt19937_64& GetRandomEngine64()
	{
		static thread_local std::mt19937_64 gen(std::random_device{}());
		return gen;
	}

	int32 RandUtil::Int32R(int32 min, int32 max)
	{
		if (min >= max)
			return min;

		auto& engine = GetRandomEngine32();
		std::uniform_int_distribution<int32> dis(min, max - 1);
		return dis(engine);
	}

	int64 RandUtil::Int64R(int64 min, int64 max)
	{
		if (min >= max)
			return min;

		auto& engine = GetRandomEngine64();
		std::uniform_int_distribution<int64> dis(min, max - 1);
		return dis(engine);
	}

	uint32 RandUtil::Uint32R(uint32 min, uint32 max)
	{
		if (min >= max)
			return min;

		auto& engine = GetRandomEngine32();
		std::uniform_int_distribution<uint32> dis(min, max - 1);
		return dis(engine);
	}

	uint64 RandUtil::Uint64R(uint64 min, uint64 max)
	{
		if (min >= max)
			return min;

		auto& engine = GetRandomEngine64();
		std::uniform_int_distribution<uint64> dis(min, max - 1);
		return dis(engine);
	}

	int32 RandUtil::Int32N(int32 max)
	{
		if (max <= 0)
			return 0;

		auto& engine = GetRandomEngine32();
		std::uniform_int_distribution<int32> dis(0, max - 1);
		return dis(engine);
	}

	int64 RandUtil::Int64N(int64 max)
	{
		if (max <= 0)
			return 0;

		auto& engine = GetRandomEngine64();
		std::uniform_int_distribution<int64> dis(0, max - 1);
		return dis(engine);
	}

	uint32 RandUtil::Uint32N(uint32 max)
	{
		if (max <= 0)
			return 0;

		auto& engine = GetRandomEngine32();
		std::uniform_int_distribution<uint32> dis(0, max - 1);
		return dis(engine);
	}

	uint64 RandUtil::Uint64N(uint64 max)
	{
		if (max <= 0)
			return 0;

		auto& engine = GetRandomEngine64();
		std::uniform_int_distribution<uint64> dis(0, max - 1);
		return dis(engine);
	}
}
