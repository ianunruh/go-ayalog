package ayalog

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const numDefaultLogFields = 6

func ParseRecord(buffer *bytes.Buffer) (Record, error) {
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

		record.Message += msg

		// reset last display hint
		lastHint = DefaultHint
	}

	return record, nil
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
