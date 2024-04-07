namespace Tyto.Utils;

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
	public static uint16 MakeUint16(uint8 high, uint8 low)
	{
		uint16 h = high;
		uint16 l = low;
		return (uint16)((h << 8) | l);
	}

	/// <inheritdoc cref="MakeUint16"/>
	public static uint32 MakeUint32(uint16 high, uint16 low)
	{
		uint32 h = high;
		uint32 l = low;
		return (h << 16) | l;
	}

	/// <inheritdoc cref="MakeUint16"/>
	public static uint64 MakeUint64(uint32 high, uint32 low)
	{
		uint64 h = high;
		uint64 l = low;
		return (h << 32) | l;
	}

	/// <summary>
	/// 获取整数的高位部分
	/// </summary>
	public static uint8 HighPartUint16(uint16 val)
	{
		return (uint8)(val >> 8);
	}

	/// <inheritdoc cref="HighPartUint16"/>
	public static uint16 HighPartUint32(uint32 val)
	{
		return (uint16)(val >> 16);
	}

	/// <inheritdoc cref="HighPartUint16"/>
	public static uint32 HighPartUint64(uint64 val)
	{
		return (uint32)(val >> 32);
	}

	/// <summary>
	/// 获取整数的低位部分
	/// </summary>
	public static uint8 LowPartUint16(uint16 val)
	{
		return (uint8)(val & 0xFF);
	}

	/// <inheritdoc cref="LowPartUint16"/>
	public static uint16 LowPartUint32(uint32 val)
	{
		return (uint16)(val & 0xFFFF);
	}

	/// <inheritdoc cref="LowPartUint16"/>
	public static uint32 LowPartUint64(uint64 val)
	{
		return (uint32)(val & 0xFFFFffff);
	}
}