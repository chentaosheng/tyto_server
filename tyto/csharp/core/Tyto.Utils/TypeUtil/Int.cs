namespace Tyto.Utils;

using int8 = sbyte;
using int16 = short;
using int32 = int;
using int64 = long;
using uint8 = byte;
using uint16 = ushort;
using uint32 = uint;
using uint64 = ulong;

/// <summary>
/// 整数合并和分拆处理类
/// </summary>
public static partial class TypeUtil
{
	/// <summary>
	/// 将high和low组合为一个新的整数，high放高位部分，low放低位部分
	/// </summary>
	public static int16 MakeInt16(int8 high, int8 low)
	{
		uint16 h = (uint8)high;
		uint16 l = (uint8)low;
		return (int16)((h << 8) | l);
	}

	/// <inheritdoc cref="MakeInt16"/>
	public static int32 MakeInt32(int16 high, int16 low)
	{
		uint32 h = (uint16)high;
		uint32 l = (uint16)low;
		return (int32)((h << 16) | l);
	}

	/// <inheritdoc cref="MakeInt16"/>
	public static int64 MakeInt64(int32 high, int32 low)
	{
		uint64 h = (uint32)high;
		uint64 l = (uint32)low;
		return (int64)((h << 32) | l);
	}

	/// <summary>
	/// 获取整数的高位部分
	/// </summary>
	public static int8 HighPartInt16(int16 val)
	{
		return (int8)((uint16)val >> 8);
	}

	/// <inheritdoc cref="HighPartInt16"/>
	public static int16 HighPartInt32(int32 val)
	{
		return (int16)((uint32)val >> 16);
	}

	/// <inheritdoc cref="HighPartInt16"/>
	public static int32 HighPartInt64(int64 val)
	{
		return (int32)((uint64)val >> 32);
	}

	/// <summary>
	/// 获取整数的低位部分
	/// </summary>
	public static int8 LowPartInt16(int16 val)
	{
		return (int8)((uint16)val & 0xFF);
	}

	/// <inheritdoc cref="LowPartInt16"/>
	public static int16 LowPartInt32(int32 val)
	{
		return (int16)((uint32)val & 0xFFFF);
	}

	/// <inheritdoc cref="LowPartInt16"/>
	public static int32 LowPartInt64(int64 val)
	{
		return (int32)((uint64)val & 0xFFFFffff);
	}
}
