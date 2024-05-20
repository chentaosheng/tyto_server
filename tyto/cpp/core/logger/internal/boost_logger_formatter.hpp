#pragma once

#include <boost/date_time/posix_time/posix_time_types.hpp>
#include <boost/log/core/record.hpp>
#include <boost/log/utility/formatting_ostream.hpp>
#include <boost/log/expressions.hpp>
#include "logger_time.hpp"

namespace tyto::internal
{
	// 格式化日志
	template<typename TLevel>
	void RecordFormatter(boost::log::record_view const& rec, boost::log::formatting_ostream& strm)
	{
		// 等价于下面语句：
		// expr::stream
		//	<< expr::format_date_time<boost::posix_time::ptime>("TimeStamp", "%Y-%m-%d %H:%M:%S")
		//	<< " [" << expr::attr<LogLevel>("Severity") << "] "
		//	<< expr::smessage;

		namespace logging = boost::log;
		namespace expr = boost::log::expressions;

		char buff[32];
		auto time_ref = logging::extract<boost::posix_time::ptime>("TimeStamp", rec);
		FormatTime(buff, sizeof(buff), time_ref.get_ptr());

		// 写入时间
		strm << buff;

		// 写入日志级别
		strm << " [" << logging::extract<TLevel>("Severity", rec) << "] ";

		// 日志内容
		strm << rec[expr::smessage];
	}
}
