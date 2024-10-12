package ayalog

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

	ArrU8Len4Arg
	ArrU8Len6Arg
	ArrU8Len16Arg
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
