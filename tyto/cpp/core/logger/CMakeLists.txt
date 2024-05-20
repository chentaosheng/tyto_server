project("tyto_logger")

add_library(tyto_logger STATIC
		internal/boost_logger.hpp
		internal/boost_logger_formatter.hpp
		internal/logger_time.cpp
		internal/logger_time.hpp
		level.cpp
		level.hpp
		logger.hpp
)

# 只依赖boost和标准库
