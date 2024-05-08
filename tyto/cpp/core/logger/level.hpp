#pragma once

#include <ostream>

namespace tyto
{
	// 日志级别
	enum LogLevel
	{
		LOG_LEVEL_TRACE = 0,
		LOG_LEVEL_DEBUG = 1,
		LOG_LEVEL_INFO = 2,
		LOG_LEVEL_WARN = 3,
		LOG_LEVEL_ERROR = 4,
		LOG_LEVEL_FATAL = 5,
	};

	// 日志级别范围
	constexpr LogLevel kMinLogLevel = LOG_LEVEL_TRACE;
	constexpr LogLevel kMaxLogLevel = LOG_LEVEL_FATAL;
	constexpr int kInvalidLogLevel = kMaxLogLevel + 1;

	// 格式化输出
	std::ostream& operator<<(std::ostream& strm, LogLevel level);
}
