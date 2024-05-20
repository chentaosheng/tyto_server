#pragma once

#include <string>
#include <format>
#include <boost/noncopyable.hpp>
#include "internal/boost_logger.hpp"
#include "level.hpp"

namespace tyto
{
	// logger包装类，支持对不同的logger进行封装
	template<typename TLogger, typename TLevel>
	class LoggerFacade : public boost::noncopyable
	{
	public:
		LoggerFacade() = default;
		~LoggerFacade() = default;

		// 初始化
		bool Init(const std::string& channel_name, const std::string& out_dir, const std::string& log_file,
			TLevel level, TLevel error_level)
		{
			return impl_.Init(channel_name, out_dir, log_file, level, error_level);
		}

		// 关闭logger
		void Close()
		{
			impl_.Close();
		}

		// 设置日志级别
		void SetLevel(TLevel level)
		{
			impl_.SetLevel(level);
		}

		// log
		void Log(TLevel level, const std::string& message)
		{
			impl_.Log(level, message);
		}

		// log with format
		template<typename... Args>
		void Logf(TLevel level, const std::format_string<Args...> fmt, Args&& ... args)
		{
			impl_.Log(level, std::format(fmt, std::forward<Args>(args)...));
		}

		// 获取实现
		TLogger& GetImpl() noexcept
		{
			return impl_;
		}

	private:
		// 实际实现
		TLogger impl_;
	};

	// 包装boost::log实现的logger
	using Logger = LoggerFacade<internal::BoostLogger<LogLevel>, LogLevel>;
}

// 对外接口
#define LOG_TRACE(logger) BOOST_LOG_IMPL(logger.GetImpl(), tyto::LOG_LEVEL_TRACE)
#define LOG_DEBUG(logger) BOOST_LOG_IMPL(logger.GetImpl(), tyto::LOG_LEVEL_DEBUG)
#define LOG_INFO(logger) BOOST_LOG_IMPL(logger.GetImpl(), tyto::LOG_LEVEL_INFO)
#define LOG_WARN(logger) BOOST_LOG_IMPL_WITH_FUNC(logger.GetImpl(), tyto::LOG_LEVEL_WARN)
#define LOG_ERROR(logger) BOOST_LOG_IMPL_WITH_FUNC(logger.GetImpl(), tyto::LOG_LEVEL_ERROR)
#define LOG_FATAL(logger) BOOST_LOG_IMPL_WITH_FUNC(logger.GetImpl(), tyto::LOG_LEVEL_FATAL)
