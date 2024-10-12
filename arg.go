package ayalog

// https://github.com/aya-rs/aya/blob/aya-v0.13.0/aya-log-common/src/lib.rs

type Field uint8

const (
	TargetField Field = iota + 1
	LevelField
	ModuleField
	FileField
	LineField
	NumArgsField
)

type Arg uint8

const (
	DisplayHintArg Arg = iota

	I8Arg
	I16Arg
	I32Arg
	I64Arg
	IsizeArg

	U8Arg
	U16Arg
	U32Arg
	U64Arg
	UsizeArg

	F32Arg
	F64Arg

	IPv4Addr
	IPv6Addr

	// `[u8; 4]` array which represents an IPv4 address.
	ArrU8Len4Arg
	// `[u8; 6]` array which represents a MAC address.
	ArrU8Len6Arg
	// `[u8; 16]` array which represents an IPv6 address.
	ArrU8Len16Arg
	// `[u16; 8]` array which represents an IPv6 address.
	ArrU16Len8Arg

	BytesArg
	StrArg
)

type DisplayHint uint8

const (
	DefaultHint DisplayHint = iota + 1
	LowerHexHint
	UpperHexHint
	IPHint
	LowerMACHint
	UpperMACHint
)
