#pragma once

#include <mutex>
#include <boost/noncopyable.hpp>

//
// 单例模板类
namespace tyto
{
	/*
		使用方法:

		class Mgr : public tyto::Singleton<Mgr>
		{
			friend class tyto::Singleton<Mgr>;
		public:
			void F() {}

		private:
			Mgr() = default;
		};
	*/
	template<typename T>
	class Singleton : public boost::noncopyable
	{
	public:
		static T& Instance()
		{
			if (instance_ == nullptr)
			{
				std::call_once(flag_, []() {
					instance_ = new T;
				});
			}

			return *instance_;
		}

	protected:
		Singleton() = default;

		virtual ~Singleton()
		{
			if (instance_ != nullptr)
			{
				delete instance_;
				instance_ = nullptr;
			}
		}

	private:
		inline static T* instance_{ nullptr };
		inline static std::once_flag flag_;
	};
}
