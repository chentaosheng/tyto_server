namespace Tyto.Utils;

using System;

public static class OsUtil
{
	/// <summary>
	/// 获取操作系统的位数
	/// </summary>
	public static int GetOsBit()
	{
		return IntPtr.Size * 8;
	}

	/// <summary>
	/// 是否为64位操作系统
	/// </summary>
	public static bool Is64Bit()
	{
		return GetOsBit() == 64;
	}

	/// <summary>
	/// 是否为32位操作系统
	/// </summary>
	public static bool Is32Bit()
	{
		return GetOsBit() == 32;
	}

	/// <summary>
	/// 是否为Windows操作系统
	/// </summary>
	public static bool IsWindows()
	{
		return Environment.OSVersion.Platform == PlatformID.Win32NT;
	}

	/// <summary>
	/// 是否为Linux操作系统
	/// </summary>
	public static bool IsLinux()
	{
		return Environment.OSVersion.Platform == PlatformID.Unix;
	}

	/// <summary>
	/// 是否为MacOSX操作系统
	/// </summary>
	public static bool IsMacOsx()
	{
		return Environment.OSVersion.Platform == PlatformID.MacOSX;
	}
}