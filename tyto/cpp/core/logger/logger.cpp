#include <cassert>
#include <functional>
#include <filesystem>
#include <regex>
#include <boost/predef/os.h>
#include <boost/date_time/posix_time/posix_time_types.hpp>
#include <boost/log/core.hpp>
#include <boost/log/expressions.hpp>
#include <boost/log/support/date_time.hpp>
#include <boost/log/utility/setup/common_attributes.hpp>
#include <boost/log/utility/setup/console.hpp>
#include <boost/algorithm/string/replace.hpp>
#include "logger.hpp"

namespace logging = boost::log;
namespace expr = boost::log::expressions;
namespace keywords = boost::log::keywords;
namespace sinks = boost::log::sinks;
namespace attrs = boost::log::attributes;

namespace tyto
{
	static void FormatTime(char* buf, std::size_t size, const boost::posix_time::ptime* ptime_ptr)
	{
		auto date = ptime_ptr->date();
		auto time = ptime_ptr->time_of_day();

		std::int32_t year = date.year();
		std::int32_t month = date.month();
		std::int32_t day = date.day();
		std::int32_t hours = static_cast<std::int32_t>(time.hours());
		std::int32_t minutes = static_cast<std::int32_t>(time.minutes());
		std::int32_t seconds = static_cast<std::int32_t>(time.seconds());
		std::int32_t milli = static_cast<std::int32_t>(time.fractional_seconds() / 1000);

		assert(size >= 24); // YYYY-MM-DD HH:MM:SS.mmm

		std::snprintf(buf, size, "%04d-%02d-%02d %02d:%02d:%02d.%03d", year, month, day, hours, minutes, seconds, milli);
	}

	static void RecordFormatter(logging::record_view const& rec, logging::formatting_ostream& strm)
	{
		// 等价于：
		// expr::stream
		//	<< expr::format_date_time<boost::posix_time::ptime>("TimeStamp", "%Y-%m-%d %H:%M:%S")
		//	<< " [" << expr::attr<LogLevel>("Severity") << "] "
		//	<< expr::smessage;

		char buf[32];
		auto time_ref = logging::extract<boost::posix_time::ptime>("TimeStamp", rec);
		FormatTime(buf, sizeof(buf), time_ref.get_ptr());

		// 写入时间
		strm << buf;

		// 写入日志级别
		strm << " [" << logging::extract<LogLevel>("Severity", rec) << "] ";

		// 日志内容
		strm << rec[expr::smessage];
	}

	bool Logger::Init(const std::string& channel_name, const std::string& out_dir, const std::string& log_file, LogLevel level)
	{
		std::scoped_lock lock(mutex_);

		if (inited_)
			return false;

		// 文件名不能包含路径
		if (log_file.find('\\') != std::string::npos || log_file.find('/') != std::string::npos)
			return false;

		log_file_ = log_file;

		std::filesystem::path dir = std::filesystem::absolute(out_dir);
		dir.make_preferred();
		out_dir_ = dir.string();

		channel_name_ = channel_name;

		// source
		impl_ = std::make_unique<LogSource>(keywords::channel = channel_name);
		impl_->add_attribute("TimeStamp", attrs::local_clock());

		// normal
		auto sink = CreateFileSink(channel_name, log_file, level, false);
		logging::core::get()->add_sink(sink);
		normal_sink_ = sink;

		// error
		sink = CreateFileSink(channel_name, "err_" + log_file, LOG_LEVEL_WARN, true);
		logging::core::get()->add_sink(sink);
		error_sink_ = sink;

#if BOOST_OS_WINDOWS
		// console
		auto streamSink = CreateStreamSink(std::cout, level);
		logging::core::get()->add_sink(streamSink);
		console_sink_ = streamSink;
#endif

		inited_ = true;

		return true;
	}

	void Logger::Close()
	{
		std::scoped_lock lock(mutex_);

		if (!inited_)
			return;

		// normal
		if (normal_sink_ != nullptr)
		{
			RemoveFileSink(normal_sink_);
			normal_sink_.reset();
		}

		// error
		if (error_sink_ != nullptr)
		{
			RemoveFileSink(error_sink_);
			error_sink_.reset();
		}

		// console
		if (console_sink_ != nullptr)
		{
			RemoveStreamSink(console_sink_);
			console_sink_.reset();
		}

		inited_ = false;
	}

	void Logger::SetLevel(LogLevel level)
	{
		std::scoped_lock lock(mutex_);

		if (!inited_)
			return;

		// normal
		if (normal_sink_ != nullptr)
		{
			normal_sink_->set_filter(expr::attr<LogLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name_);
		}

		// error
		if (error_sink_ != nullptr)
		{
			if (level > LOG_LEVEL_WARN)
			{
				error_sink_->set_filter(expr::attr<LogLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name_);
			}
			else
			{
				error_sink_->set_filter(expr::attr<LogLevel>("Severity") >= LOG_LEVEL_WARN && expr::attr<std::string>("Channel") == channel_name_);
			}
		}

		// console
		if (console_sink_ != nullptr)
		{
			console_sink_->set_filter(expr::attr<LogLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name_);
		}
	}

	void Logger::SetMaxFileAge(std::chrono::seconds sec)
	{
		if (sec <= std::chrono::seconds(0))
			return;

		std::scoped_lock lock(mutex_);

		max_file_age_ = sec;
	}

	void Logger::CleanupFile(boost::shared_ptr<Logger::FileBackend> backend,
		const std::string& log_file, sinks::text_file_backend::stream_type& file)
	{
		std::scoped_lock lock(mutex_);

		std::filesystem::path out_path = out_dir_;
		std::string pattern = out_path.append(log_file).make_preferred().string();

		// 对路径中的特殊字符进行转义
		boost::algorithm::replace_all(pattern, "\\", "\\\\");
		boost::algorithm::replace_all(pattern, ".", "\\.");

		// 生成用于文件匹配的正则表达式
		pattern = "^" + pattern + "\\..*";
		std::regex pattern_regex(pattern);
		std::filesystem::path path_to_search(out_dir_);

		// 过期时间点
		auto expire_time = std::chrono::system_clock::now();
		expire_time -= max_file_age_;

		// 当前正在写入的文件
		auto cur_file = backend->get_current_file_name().string();

		// 搜索目录下的所有文件
		for (const auto& entry : std::filesystem::recursive_directory_iterator(path_to_search))
		{
			if (!std::regex_match(entry.path().string(), pattern_regex))
				continue;

			if (!std::filesystem::is_regular_file(entry.status()))
				continue;

			// 跳过当前文件
			if (cur_file == entry.path())
				continue;

			// 判断时间
			std::error_code err;
			auto file_time = entry.last_write_time(err);
			if (err)
				continue;

			auto mod_time = std::chrono::clock_cast<std::chrono::system_clock>(file_time);
			if (mod_time > expire_time)
				continue;

			// 删除过时文件
			std::filesystem::remove(entry.path());
		}
	}

	boost::shared_ptr<tyto::Logger::FileSink> Logger::CreateFileSink(const std::string& channel_name,
		const std::string& log_file, LogLevel level, bool auto_flush)
	{
		std::filesystem::path out_path(out_dir_);
		auto pattern = out_path.append(log_file).string();
		pattern.append(".%Y-%m-%d");

		auto backend = boost::make_shared<FileBackend>(
			keywords::file_name = pattern,
			keywords::target_file_name = pattern,
			keywords::time_based_rotation = sinks::file::rotation_at_time_point(0, 0, 0),
			keywords::enable_final_rotation = false,
			keywords::open_mode = std::ios::app | std::ios::out
		);

		if (auto_flush)
			backend->auto_flush(true);

		// 每次打开新文件后，清理过期文件
		backend->set_open_handler(std::bind(&Logger::CleanupFile, this, backend, log_file, std::placeholders::_1));

		auto sink = boost::make_shared<FileSink>(backend);
		sink->set_formatter(&RecordFormatter);

		sink->set_filter(expr::attr<LogLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name);

		return sink;
	}

	void Logger::RemoveFileSink(boost::shared_ptr<FileSink> sink)
	{
		logging::core::get()->remove_sink(sink);
		sink->stop();
		sink->flush();
	}

	boost::shared_ptr<Logger::StreamSink> Logger::CreateStreamSink(std::ostream& strm, LogLevel level)
	{
		auto strm_ptr = boost::shared_ptr<std::ostream>(&strm, boost::null_deleter());
		auto backend = boost::make_shared<StreamBackend>();

		backend->add_stream(strm_ptr);

		auto sink = boost::make_shared<StreamSink>(backend);

		sink->set_filter(expr::attr<LogLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name_);

		sink->set_formatter(&RecordFormatter);

		return sink;
	}

	void Logger::RemoveStreamSink(boost::shared_ptr<StreamSink> sink)
	{
		logging::core::get()->remove_sink(sink);
		sink->flush();
	}
}
