#pragma once

#include <ostream>

namespace tyto
{
	// 日志级别
	enum class LogLevel
	{
		TRACE = 0,
		DEBUG = 1,
		INFO = 2,
		WARN = 3,
		ERROR = 4,
		FATAL = 5,
	};

	// 日志级别范围
	constexpr LogLevel kMinLogLevel = LogLevel::TRACE;
	constexpr LogLevel kMaxLogLevel = LogLevel::FATAL;
	constexpr int kInvalidLogLevel = static_cast<int>(kMaxLogLevel) + 1;

	// 格式化输出
	std::ostream& operator<<(std::ostream& strm, LogLevel level);
}
