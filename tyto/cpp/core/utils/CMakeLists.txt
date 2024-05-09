project("tyto_utils")

add_library(tyto_utils STATIC
		typedef.hpp
		cpu_relax.hpp
		spin_lock.hpp
		spin_lock.cpp
		singleton.hpp
		time_util.hpp
		time_util.cpp
		rand_util.cpp
		rand_util.hpp
)
