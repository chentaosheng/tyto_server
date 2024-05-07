#include <boost/predef/os.h>
#include <boost/date_time/posix_time/posix_time_types.hpp>
#include <boost/log/core.hpp>
#include <boost/log/expressions.hpp>
#include <boost/log/support/date_time.hpp>
#include <boost/log/utility/setup/common_attributes.hpp>
#include <boost/log/utility/setup/console.hpp>
#include <boost/filesystem/path.hpp>
#include "logger.hpp"

namespace logging = boost::log;
namespace expr = boost::log::expressions;
namespace keywords = boost::log::keywords;
namespace sinks = boost::log::sinks;
namespace attrs = boost::log::attributes;

namespace tyto
{
	bool Logger::Init(const std::string& channel_name, const std::string& out_dir, const std::string& log_file, LogLevel level)
	{
		std::scoped_lock lock(mutex_);

		if (inited_)
			return false;

		channel_name_ = channel_name;

		// source
		impl_ = std::make_unique<LogSource>(keywords::channel = channel_name);
		impl_->add_attribute("TimeStamp", attrs::local_clock());

		// normal
		auto sink = CreateFileSink(channel_name, out_dir, log_file, level, false);
		logging::core::get()->add_sink(sink);
		normal_sink_ = sink;

		// error
		sink = CreateFileSink(channel_name, out_dir, "err_" + log_file, LogLevel::WARN, true);
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
			if (level > LogLevel::WARN)
			{
				error_sink_->set_filter(expr::attr<LogLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name_);
			}
			else
			{
				error_sink_->set_filter(expr::attr<LogLevel>("Severity") >= LogLevel::WARN && expr::attr<std::string>("Channel") == channel_name_);
			}
		}

		// console
		if (console_sink_ != nullptr)
		{
			console_sink_->set_filter(expr::attr<LogLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name_);
		}
	}

	boost::shared_ptr<tyto::Logger::FileSink> Logger::CreateFileSink(const std::string& channel_name,
		const std::string& out_dir, const std::string& log_file, LogLevel level, bool auto_flush)
	{
		boost::filesystem::path out_path(out_dir);
		auto str_path = out_path.append(log_file).make_preferred().string();

		auto backend = boost::make_shared<FileBackend>(
			keywords::file_name = str_path,
			keywords::target_file_name = str_path + ".%Y-%m-%d",
			keywords::time_based_rotation = sinks::file::rotation_at_time_point(0, 0, 0),
			keywords::enable_final_rotation = false,
			keywords::max_files = kMaxFileCount,
			keywords::open_mode = std::ios::app | std::ios::out
		);

		if (auto_flush)
		{
			backend->auto_flush(true);
		}

		auto sink = boost::make_shared<FileSink>(backend);
		sink->set_formatter(
			expr::stream
			<< expr::format_date_time<boost::posix_time::ptime>("TimeStamp", "%Y-%m-%d %H:%M:%S")
			<< " [" << expr::attr<LogLevel>("Severity") << "] "
			<< expr::smessage
		);

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

		sink->set_formatter(
			expr::stream
			<< expr::format_date_time<boost::posix_time::ptime>("TimeStamp", "%Y-%m-%d %H:%M:%S")
			<< " [" << expr::attr<LogLevel>("Severity") << "] "
			<< expr::smessage
		);

		return sink;
	}

	void Logger::RemoveStreamSink(boost::shared_ptr<StreamSink> sink)
	{
		logging::core::get()->remove_sink(sink);
		sink->flush();
	}
}
