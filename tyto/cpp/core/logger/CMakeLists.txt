project("tyto_logger")

add_library(tyto_logger STATIC
		logger.cpp
		logger.hpp
		level.cpp
		level.hpp
)

# 只依赖boost和标准库
