package ayalog

import (
	"encoding/binary"
	"fmt"
	"math"
	"net/netip"
)

// https://github.com/aya-rs/aya/blob/aya-log-v0.2.0/aya-log/src/lib.rs

func formatArg(arg Arg, hint DisplayHint, v []byte) (string, error) {
	switch arg {
	case I8Arg:
		if len(v) != 1 {
			return "", fmt.Errorf("expected i8 arg to be exactly 1 byte: %v", v)
		}
		return formatInt(hint, int8(v[0]))
	case I16Arg:
		return formatInt(hint, int16(binary.LittleEndian.Uint16(v)))
	case I32Arg:
		return formatInt(hint, int32(binary.LittleEndian.Uint32(v)))
	case I64Arg, IsizeArg:
		return formatInt(hint, int64(binary.LittleEndian.Uint64(v)))
	case U8Arg:
		if len(v) != 1 {
			return "", fmt.Errorf("expected u8 arg to be exactly 1 byte: %v", v)
		}
		return formatInt(hint, uint8(v[0]))
	case U16Arg:
		return formatInt(hint, binary.LittleEndian.Uint16(v))
	case U32Arg:
		return formatU32Arg(hint, binary.LittleEndian.Uint32(v))
	case U64Arg, UsizeArg:
		return formatInt(hint, binary.LittleEndian.Uint64(v))
	case F32Arg:
		return formatFloat(hint, math.Float32frombits(binary.LittleEndian.Uint32(v)))
	case F64Arg:
		return formatFloat(hint, math.Float64frombits(binary.LittleEndian.Uint64(v)))
	case ArrU8Len6Arg:
		if len(v) != 6 {
			return "", fmt.Errorf("expected ArrU8Len6 arg to be exactly 6 bytes: %v", v)
		}
		return formatArrU8Len6Arg(hint, [6]byte(v))
	case ArrU8Len16Arg:
		if len(v) != 16 {
			return "", fmt.Errorf("expected ArrU8Len16 arg to be exactly 16 bytes: %v", v)
		}
		return formatArrU8Len16Arg(hint, [16]byte(v))
	// TODO implement me
	// case ArrU16Len8Arg:
	case BytesArg:
		return formatBytes(hint, v)
	case StrArg:
		return string(v), nil
	default:
		return "", fmt.Errorf("unknown arg tag: %d", arg)
	}
}

func formatU32Arg(hint DisplayHint, v uint32) (string, error) {
	switch hint {
	case IPHint:
		return formatIPv4Addr(v), nil
	}
	return formatInt(hint, v)
}

func formatArrU8Len6Arg(hint DisplayHint, v [6]byte) (string, error) {
	switch hint {
	case LowerMACHint:
		return formatLowerMAC(v), nil
	case UpperMACHint:
		return formatUpperMAC(v), nil
	}
	return "", fmt.Errorf("unsupported display hint: %d", hint)
}

func formatArrU8Len16Arg(hint DisplayHint, v [16]byte) (string, error) {
	switch hint {
	case IPHint:
		return netip.AddrFrom16(v).String(), nil
	}
	return "", fmt.Errorf("unsupported display hint: %d", hint)
}

func formatFloat(hint DisplayHint, v any) (string, error) {
	switch hint {
	case DefaultHint:
		return formatDefault(v), nil
	}
	return "", fmt.Errorf("unsupported display hint: %d", hint)
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

func formatBytes(hint DisplayHint, v []byte) (string, error) {
	switch hint {
	case LowerHexHint:
		return formatLowerHex(v), nil
	case UpperHexHint:
		return formatUpperHex(v), nil
	}
	return "", fmt.Errorf("unsupported display hint: %d", hint)
}

func formatIPv4Addr(v uint32) string {
	b := [4]byte{}
	binary.BigEndian.PutUint32(b[:], v)
	return netip.AddrFrom4(b).String()
}

func formatLowerMAC(v [6]byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		v[0], v[1], v[2], v[3], v[4], v[5])
}

func formatUpperMAC(v [6]byte) string {
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
		v[0], v[1], v[2], v[3], v[4], v[5])
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
