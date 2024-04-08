namespace Tyto.Utils.PatternUtil;

using System;

/// <summary>
/// 单例模式，派生类必须将构造函数设置为protected或private
/// </summary>
public class Singleton<T> where T : class
{
	public static T Instance => SingletonCreator.s_instance;

	protected Singleton()
	{
	}

	private static class SingletonCreator
	{
		internal static readonly T s_instance = (T)Activator.CreateInstance(typeof(T), true)!;

		static SingletonCreator()
		{
		}
	}
}