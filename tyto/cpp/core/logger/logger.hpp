#pragma once

#include <cstdint>
#include <string>
#include <memory>
#include <mutex>
#include <chrono>
#include <format>
#include <boost/predef/compiler.h>
#include <boost/noncopyable.hpp>
#include <boost/shared_ptr.hpp>
#include <boost/log/sinks.hpp>
#include <boost/log/sources/severity_channel_logger.hpp>
#include <boost/log/sources/record_ostream.hpp>
#include "level.hpp"

namespace tyto
{
	// 日志类
	class Logger : public boost::noncopyable
	{
	public:
		using LogSource = boost::log::sources::severity_channel_logger_mt<LogLevel, std::string>;
		using FileBackend = boost::log::sinks::text_file_backend;
		using FileSink = boost::log::sinks::asynchronous_sink<FileBackend>;
		using StreamBackend = boost::log::sinks::text_ostream_backend;
		using StreamSink = boost::log::sinks::synchronous_sink<StreamBackend>;

		// 默认最大保留时间，单位：秒
		static constexpr std::int64_t kDefaultMaxFileAge = 30 * 24 * 3600;

	public:
		Logger() = default;
		~Logger() = default;

		// 初始化
		bool Init(const std::string& channel_name, const std::string& out_dir, const std::string& log_file, LogLevel level);
		// 关闭logger
		void Close();

		// 设置日志级别
		void SetLevel(LogLevel level);
		// 设置文件最大保留时间，单位：秒
		void SetMaxFileAge(std::chrono::seconds sec);

		inline LogSource& Source() noexcept
		{
			return *impl_.get();
		}

		inline const std::string& ChannelName() const noexcept
		{
			return channel_name_;
		}

		// 打印日志
		void Log(LogLevel level, const std::string& message)
		{
			BOOST_LOG_SEV(*impl_.get(), level) << message;
		}

		// C++方式格式化打印日志
		template<typename... Args>
		void Logf(LogLevel level, const std::format_string<Args...> fmt, Args&&... args)
		{
			BOOST_LOG_SEV(*impl_.get(), level) << std::vformat(fmt, std::make_format_args(args...));
		}

	private:
		void CleanupFile(const boost::shared_ptr<FileSink>& sink, const std::string& log_file,
			boost::log::sinks::text_file_backend::stream_type& file);

		boost::shared_ptr<FileSink> CreateFileSink(const std::string& channel_name,
			const std::string& log_file, LogLevel level, bool auto_flush);
		void RemoveFileSink(boost::shared_ptr<FileSink>& sink);

		boost::shared_ptr<StreamSink> CreateStreamSink(std::ostream& strm, LogLevel level);
		void RemoveStreamSink(boost::shared_ptr<StreamSink>& sink);

	private:
		bool inited_{ false };

		std::mutex mutex_;
		std::string channel_name_;
		std::string log_file_;
		std::string out_dir_;
		std::atomic<std::chrono::seconds> max_file_age_;

		std::unique_ptr<LogSource> impl_;
		boost::shared_ptr<FileSink> normal_sink_;
		boost::shared_ptr<FileSink> error_sink_;
		boost::shared_ptr<StreamSink> console_sink_;
	};
}

// 普通打印
#define LOG(logger, level) BOOST_LOG_SEV(logger.Source(), level)

// 打印附带函数名和行号
#if BOOST_COMP_MSVC || BOOST_COMP_MSVC_EMULATED
#	define LOG_FUNC(logger, level) BOOST_LOG_SEV(logger.Source(), level) << "[" << __FUNCSIG__ << ":" << __LINE__ << "] "
#elif BOOST_COMP_GNUC || BOOST_COMP_GNUC_EMULATED
#	define LOG_FUNC(logger, level) BOOST_LOG_SEV(logger.Source(), level) << "[" << __PRETTY_FUNCTION__ << ":" << __LINE__ << "] "
#else
#	error The compiler is not supported
#endif

// 对外接口
#define LOG_TRACE(logger) LOG(logger, tyto::LOG_LEVEL_TRACE)
#define LOG_DEBUG(logger) LOG(logger, tyto::LOG_LEVEL_DEBUG)
#define LOG_INFO(logger) LOG(logger, tyto::LOG_LEVEL_INFO)
#define LOG_WARN(logger) LOG_FUNC(logger, tyto::LOG_LEVEL_WARN)
#define LOG_ERROR(logger) LOG_FUNC(logger, tyto::LOG_LEVEL_ERROR)
#define LOG_FATAL(logger) LOG_FUNC(logger, tyto::LOG_LEVEL_FATAL)
