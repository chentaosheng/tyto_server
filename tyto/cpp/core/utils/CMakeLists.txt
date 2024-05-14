project("tyto_utils")

add_library(tyto_utils STATIC
		typedef.hpp
		cpu_relax.hpp
		spin_lock.cpp
		spin_lock.hpp
		singleton.hpp
		time_util.cpp
		time_util.hpp
		rand_util.cpp
		rand_util.hpp
		id_pool.cpp
		id_pool.hpp
)
