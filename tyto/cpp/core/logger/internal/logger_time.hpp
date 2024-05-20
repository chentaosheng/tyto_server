#pragma once

#include <boost/date_time/posix_time/posix_time_types.hpp>

namespace tyto::internal
{
	// 格式化时间为：YYYY-MM-DD HH:MM:SS.fff
	void FormatTime(char* buf, std::size_t size, const boost::posix_time::ptime* ptime_ptr);
}
