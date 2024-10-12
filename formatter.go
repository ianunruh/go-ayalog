package ayalog

import (
	"encoding/binary"
	"fmt"
	"net/netip"
)

// https://github.com/aya-rs/aya/blob/aya-v0.13.0/aya-log/src/lib.rs

func formatArg(arg Arg, hint DisplayHint, v []byte) (string, error) {
	switch arg {
	case I8Arg:
		return formatInt(hint, int8(v[0]))
	case I16Arg:
		return formatInt(hint, int16(binary.LittleEndian.Uint16(v)))
	case I32Arg:
		return formatInt(hint, int32(binary.LittleEndian.Uint32(v)))
	case I64Arg, IsizeArg:
		return formatInt(hint, int64(binary.LittleEndian.Uint64(v)))
	case U8Arg:
		return formatInt(hint, uint8(v[0]))
	case U16Arg:
		return formatInt(hint, binary.LittleEndian.Uint16(v))
	case U32Arg:
		return formatU32Arg(hint, binary.LittleEndian.Uint32(v))
	case U64Arg, UsizeArg:
		return formatInt(hint, binary.LittleEndian.Uint64(v))
	case F32Arg:
		// TODO implement me
	case F64Arg:
		// TODO implement me
	case IPv4Addr:
		// TODO implement me
	case IPv6Addr:
		// TODO implement me
	case ArrU8Len4Arg:
		// TODO implement me
	case ArrU8Len6Arg:
		// TODO implement me
	case ArrU8Len16Arg:
		// TODO implement me
	case ArrU16Len8Arg:
		// TODO implement me
	case BytesArg:
		// TODO support hex formats
		return string(v), nil
	case StrArg:
		return string(v), nil
	default:
		return "", fmt.Errorf("unknown arg tag: %d", arg)
	}

	return "", nil
}

func formatU32Arg(hint DisplayHint, v uint32) (string, error) {
	switch hint {
	case IPHint:
		return formatIPv4Addr(v), nil
	}
	return formatInt(hint, v)
}

func formatInt(hint DisplayHint, v any) (string, error) {
	switch hint {
	case LowerHexHint:
		return formatLowerHex(v), nil
	case UpperHexHint:
		return formatUpperHex(v), nil
	case DefaultHint:
		return formatDefault(v), nil
	}
	return "", fmt.Errorf("unsupported display hint: %d", hint)
}

func formatIPv4Addr(v uint32) string {
	b := [4]byte{}
	binary.BigEndian.PutUint32(b[:], v)
	return netip.AddrFrom4(b).String()
}

func formatLowerHex(v any) string {
	return fmt.Sprintf("%016x", v)
}

func formatUpperHex(v any) string {
	return fmt.Sprintf("%016X", v)
}

func formatDefault(v any) string {
	return fmt.Sprintf("%v", v)
}
