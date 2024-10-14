package ayalog

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const numDefaultLogFields = 6

const (
	LogLibraryVersion0_1_0 = "0.1.0"
	LogLibraryVersion0_1_1 = "0.1.1"
)

func ParseRecord(buffer *bytes.Buffer) (Record, error) {
	p := Parser{}
	return p.Record(buffer)
}

type Parser struct {
	// IncludeArgs controls whether each record has individual args populated.
	IncludeArgs bool

	// LogLibraryVersion is used to determine how to parse the records from
	// aya-log-ebpf, as the binary format differs between versions.
	LogLibraryVersion string
}

func (p Parser) Record(buffer *bytes.Buffer) (Record, error) {
	var (
		record  Record
		numArgs uint64
	)

	// parse built-in log fields
	for range numDefaultLogFields {
		tag, value, err := readField[Field](buffer)
		if err != nil {
			return record, err
		}

		switch tag {
		case TargetField:
			record.Target = string(value)
		case LevelField:
			if len(value) != 1 {
				return record, fmt.Errorf("expected level field to be exactly 1 byte: %v", value)
			}
			record.Level = Level(value[0])
		case ModuleField:
			record.Module = string(value)
		case FileField:
			record.File = string(value)
		case LineField:
			record.Line = binary.LittleEndian.Uint32(value)
		case NumArgsField:
			numArgs = binary.LittleEndian.Uint64(value)
		}
	}

	// parse variable log fields
	lastHint := DefaultHint
	for i := range numArgs {
		tag, value, err := readField[Arg](buffer)
		if err != nil {
			return record, err
		}

		tag = p.mapArgTag(tag)

		if tag == DisplayHintArg {
			if len(value) != 1 {
				return record, fmt.Errorf("expected display hint arg to be exactly 1 byte: %v", value)
			}
			lastHint = DisplayHint(value[0])
			continue
		}

		msg, err := formatArg(tag, lastHint, value)
		if err != nil {
			return record, fmt.Errorf("formatting arg %d, tag %d: %w", i, tag, err)
		}

		if p.IncludeArgs {
			record.Args = append(record.Args, RecordArg{
				Type:        tag,
				DisplayHint: lastHint,
				Value:       value,
				Formatted:   msg,
			})
		}

		record.Message += msg

		// reset last display hint
		lastHint = DefaultHint
	}

	return record, nil
}

func (p Parser) mapArgTag(arg Arg) Arg {
	if p.LogLibraryVersion == LogLibraryVersion0_1_0 {
		// aya-log-ebpf v0.1.1 introduced 3 new args after F64, which changed the enum
		// order for the subsequent args.
		if arg > F64Arg {
			arg += 3
		}
	}
	return arg
}

func readField[T ~uint8](buf *bytes.Buffer) (T, []byte, error) {
	var tag T
	if err := binary.Read(buf, binary.LittleEndian, &tag); err != nil {
		return 0, nil, fmt.Errorf("reading tag: %w", err)
	}

	var valueLen uint16
	if err := binary.Read(buf, binary.LittleEndian, &valueLen); err != nil {
		return 0, nil, fmt.Errorf("reading log value len: %w", err)
	}

	// TODO check read bytes matches valueLen? EOF error might cover this
	value := make([]byte, valueLen)
	if _, err := buf.Read(value); err != nil {
		return 0, nil, fmt.Errorf("reading value: %w", err)
	}

	return tag, value, nil
}
