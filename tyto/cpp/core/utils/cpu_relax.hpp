#pragma once

#include <boost/predef.h>

// PAUSE指令实现
#if BOOST_COMP_MSVC || BOOST_COMP_MSVC_EMULATED
	#include <windows.h>
	#define CpuRelax() YieldProcessor()
#elif BOOST_COMP_GNUC || BOOST_COMP_GNUC_EMULATED
	#include <immintrin.h>
	#define CpuRelax() _mm_pause()
#else
	#error The compiler is not supported
#endif
