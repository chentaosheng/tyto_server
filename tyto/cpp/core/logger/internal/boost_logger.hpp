#pragma once

#include <cstdint>
#include <string>
#include <memory>
#include <mutex>
#include <chrono>
#include <format>
#include <functional>
#include <filesystem>
#include <regex>
#include <boost/predef/os.h>
#include <boost/predef/compiler.h>
#include <boost/noncopyable.hpp>
#include <boost/shared_ptr.hpp>
#include <boost/log/sinks.hpp>
#include <boost/log/sources/severity_channel_logger.hpp>
#include <boost/log/sources/record_ostream.hpp>
#include <boost/log/core.hpp>
#include <boost/log/expressions.hpp>
#include <boost/log/support/date_time.hpp>
#include <boost/log/utility/setup/common_attributes.hpp>
#include <boost/log/utility/setup/console.hpp>
#include <boost/algorithm/string/replace.hpp>
#include "boost_logger_formatter.hpp"

namespace tyto::internal
{
	// 日志类
	template<typename TLevel>
	class BoostLogger : public boost::noncopyable
	{
	public:
		using LogSource = boost::log::sources::severity_channel_logger_mt<TLevel, std::string>;
		using FileBackend = boost::log::sinks::text_file_backend;
		using FileSink = boost::log::sinks::asynchronous_sink<FileBackend>;
		using StreamBackend = boost::log::sinks::text_ostream_backend;
		using StreamSink = boost::log::sinks::synchronous_sink<StreamBackend>;

		// 默认最大保留时间，单位：秒
		static constexpr std::int64_t kDefaultMaxFileAge = 30 * 24 * 3600;

	public:
		BoostLogger() = default;
		~BoostLogger() = default;

		// 初始化
		bool Init(const std::string& channel_name, const std::string& out_dir, const std::string& log_file,
			TLevel level, TLevel error_level);

		// 关闭logger
		void Close();

		// 设置日志级别
		void SetLevel(TLevel level);
		// 设置文件最大保留时间，单位：秒
		void SetMaxFileAge(std::chrono::seconds sec);

		LogSource& GetSource() noexcept
		{
			return *impl_.get();
		}

		const std::string& GetChannelName() const noexcept
		{
			return channel_name_;
		}

		// 打印日志
		void Log(TLevel level, const std::string& message)
		{
			BOOST_LOG_SEV(*impl_.get(), level) << message;
		}

		// C++方式格式化打印日志
		template<typename... Args>
		void Logf(TLevel level, const std::format_string<Args...> fmt, Args&& ... args)
		{
			BOOST_LOG_SEV(*impl_.get(), level) << std::vformat(fmt, std::make_format_args(args...));
		}

	private:
		void CleanupFile(const boost::shared_ptr<FileSink>& sink, const std::string& log_file,
			boost::log::sinks::text_file_backend::stream_type& file);

		boost::shared_ptr<FileSink> CreateFileSink(const std::string& channel_name,
			const std::string& log_file, TLevel level, bool auto_flush);
		void RemoveFileSink(boost::shared_ptr<FileSink>& sink);

		boost::shared_ptr<StreamSink> CreateStreamSink(std::ostream& strm, TLevel level);
		void RemoveStreamSink(boost::shared_ptr<StreamSink>& sink);

	private:
		bool inited_{ false };

		std::mutex mutex_;
		std::string channel_name_;
		std::string log_file_;
		std::string out_dir_;
		TLevel error_level_;
		std::atomic<std::chrono::seconds> max_file_age_{ std::chrono::seconds(kDefaultMaxFileAge) };

		std::unique_ptr<LogSource> impl_;
		boost::shared_ptr<FileSink> normal_sink_;
		boost::shared_ptr<FileSink> error_sink_;
		boost::shared_ptr<StreamSink> console_sink_;
	};

	// 初始化
	template<typename TLevel>
	bool BoostLogger<TLevel>::Init(const std::string& channel_name, const std::string& out_dir, const std::string& log_file,
		TLevel level, TLevel error_level)
	{
		namespace logging = boost::log;
		namespace keywords = boost::log::keywords;
		namespace attrs = boost::log::attributes;

		std::scoped_lock lock(mutex_);

		if (inited_)
			return false;

		// 文件名不能包含路径
		if (log_file.find('\\') != std::string::npos || log_file.find('/') != std::string::npos)
			return false;

		std::filesystem::path dir = std::filesystem::absolute(out_dir);
		out_dir_ = dir.make_preferred().string();

		log_file_ = log_file;
		channel_name_ = channel_name;
		error_level_ = error_level;

		// 创建目录
		std::error_code err;
		std::filesystem::create_directories(dir, err);
		if (err)
			return false;

		// source
		impl_ = std::make_unique<LogSource>(keywords::channel = channel_name);
		impl_->add_attribute("TimeStamp", attrs::local_clock());

		// normal
		auto sink = CreateFileSink(channel_name, log_file, level, true);
		logging::core::get()->add_sink(sink);
		normal_sink_ = sink;

		// error
		sink = CreateFileSink(channel_name, "err_" + log_file, error_level_, true);
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

	// 关闭logger
	template<typename TLevel>
	void BoostLogger<TLevel>::Close()
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

	// 设置日志级别
	template<typename TLevel>
	void BoostLogger<TLevel>::SetLevel(TLevel level)
	{
		namespace expr = boost::log::expressions;

		std::scoped_lock lock(mutex_);

		if (!inited_)
			return;

		// normal
		if (normal_sink_ != nullptr)
		{
			normal_sink_->set_filter(expr::attr<TLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name_);
		}

		// error
		if (error_sink_ != nullptr)
		{
			if (level > error_level_)
			{
				error_sink_->set_filter(expr::attr<TLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name_);
			}
			else
			{
				error_sink_->set_filter(expr::attr<TLevel>("Severity") >= error_level_ && expr::attr<std::string>("Channel") == channel_name_);
			}
		}

		// console
		if (console_sink_ != nullptr)
		{
			console_sink_->set_filter(expr::attr<TLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name_);
		}
	}

	// 设置日志文件最大保留时间
	template<typename TLevel>
	void BoostLogger<TLevel>::SetMaxFileAge(std::chrono::seconds sec)
	{
		if (sec <= std::chrono::seconds(0))
			return;

		max_file_age_.store(sec, std::memory_order_relaxed);
	}

	// 清理过期日志
	template<typename TLevel>
	void BoostLogger<TLevel>::CleanupFile(const boost::shared_ptr<FileSink>& sink, const std::string& log_file, boost::log::sinks::text_file_backend::stream_type& file)
	{
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
		expire_time -= max_file_age_.load(std::memory_order_relaxed);

		// 当前正在写入的文件
		auto cur_file = sink->locked_backend()->get_current_file_name().string();

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
			auto write_time = entry.last_write_time(err);
			if (err)
				continue;

			auto file_time = std::chrono::clock_cast<std::chrono::system_clock>(write_time);
			if (file_time > expire_time)
				continue;

			// 删除过时文件
			std::error_code serr;
			std::filesystem::remove(entry.path(), serr);
		}

	}

	// 创建文件sink
	template<typename TLevel>
	boost::shared_ptr<typename BoostLogger<TLevel>::FileSink> BoostLogger<TLevel>::CreateFileSink(const std::string& channel_name, const std::string& log_file, TLevel level, bool auto_flush)
	{
		namespace expr = boost::log::expressions;
		namespace keywords = boost::log::keywords;
		namespace sinks = boost::log::sinks;

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

		auto sink = boost::make_shared<FileSink>(backend);
		sink->set_formatter(&RecordFormatter<TLevel>);
		sink->set_filter(expr::attr<TLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name);

		// 每次打开新文件后，清理过期文件
		sink->locked_backend()->set_open_handler(std::bind(&BoostLogger::CleanupFile, this, sink, log_file, std::placeholders::_1));

		return sink;
	}

	// 删除文件sink
	template<typename TLevel>
	void BoostLogger<TLevel>::RemoveFileSink(boost::shared_ptr<FileSink>& sink)
	{
		namespace logging = boost::log;

		logging::core::get()->remove_sink(sink);
		sink->stop();
		sink->flush();
	}

	// 创建流式sink，主要用于控制台输出
	template<typename TLevel>
	boost::shared_ptr<typename BoostLogger<TLevel>::StreamSink> BoostLogger<TLevel>::CreateStreamSink(std::ostream& strm, TLevel level)
	{
		namespace expr = boost::log::expressions;

		auto strm_ptr = boost::shared_ptr<std::ostream>(&strm, boost::null_deleter());
		auto backend = boost::make_shared<StreamBackend>();

		backend->add_stream(strm_ptr);

		auto sink = boost::make_shared<StreamSink>(backend);

		sink->set_filter(expr::attr<TLevel>("Severity") >= level && expr::attr<std::string>("Channel") == channel_name_);

		sink->set_formatter(&RecordFormatter<TLevel>);

		return sink;
	}

	// 删除流式sink
	template<typename TLevel>
	void BoostLogger<TLevel>::RemoveStreamSink(boost::shared_ptr<StreamSink>& sink)
	{
		namespace logging = boost::log;

		logging::core::get()->remove_sink(sink);
		sink->flush();
	}
}

// 根据不同编译器选择不同的函数名宏
#if BOOST_COMP_MSVC || BOOST_COMP_MSVC_EMULATED
#	define FULL_FUNCTION_NAME __FUNCSIG__
#elif BOOST_COMP_GNUC || BOOST_COMP_GNUC_EMULATED
#	define FULL_FUNCTION_NAME __PRETTY_FUNCTION__
#else
#	error The compiler is not supported
#endif

// 普通打印
#define BOOST_LOG_IMPL(logger, level) BOOST_LOG_SEV(logger.GetSource(), level)

// 打印附带函数名和行号
#define BOOST_LOG_IMPL_WITH_FUNC(logger, level) BOOST_LOG_SEV(logger.GetSource(), level) << "[" << FULL_FUNCTION_NAME << ":" << __LINE__ << "] "
