project("tyto_logger")

add_library(tyto_logger STATIC
		logger.hpp
		logger.cpp
		level.hpp
		level.cpp
)

# 只依赖boost和标准库
