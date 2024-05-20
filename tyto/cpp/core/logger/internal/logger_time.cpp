#include <cstdio>
#include <boost/assert.hpp>
#include <boost/date_time/posix_time/posix_time.hpp>
#include "logger_time.hpp"

namespace tyto::internal
{
	void FormatTime(char* buf, std::size_t size, const boost::posix_time::ptime* ptime_ptr)
	{
		auto date = ptime_ptr->date();
		auto time = ptime_ptr->time_of_day();

		auto year = static_cast<std::int32_t>(date.year());
		auto month = static_cast<std::int32_t>(date.month());
		auto day = static_cast<std::int32_t>(date.day());
		auto hours = static_cast<std::int32_t>(time.hours());
		auto minutes = static_cast<std::int32_t>(time.minutes());
		auto seconds = static_cast<std::int32_t>(time.seconds());
		auto milli = static_cast<std::int32_t>(time.fractional_seconds() / 1000);

		// YYYY-MM-DD HH:MM:SS.mmm
		BOOST_ASSERT(size >= 24);

		std::snprintf(buf, size, "%04d-%02d-%02d %02d:%02d:%02d.%03d", year, month, day, hours, minutes, seconds, milli);
	}
}
