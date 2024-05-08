#include "level.hpp"

namespace tyto
{
	std::ostream& operator<<(std::ostream& strm, LogLevel level)
	{
		static const char* names[] = {
			"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "UNKNOWN",
		};

		if (level < kMinLogLevel || level > kMaxLogLevel)
			strm << names[static_cast<int>(kInvalidLogLevel)];
		else
			strm << names[static_cast<int>(level)];

		return strm;
	}
}
