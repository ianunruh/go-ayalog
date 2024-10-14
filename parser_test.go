package ayalog

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO add test data for v0.1.1

func TestParseRecord_0_1_0(t *testing.T) {
	record, err := os.ReadFile("testdata/0_1_0.bin")
	require.NoError(t, err)

	parser := Parser{
		IncludeArgs:       true,
		LogLibraryVersion: LogLibraryVersion0_1_0,
	}

	r, err := parser.Record(bytes.NewBuffer(record))
	require.NoError(t, err)

	assert.Equal(t, "xdp_hello", r.Target)
	assert.Equal(t, "xdp_hello", r.Module)
	assert.Equal(t, "src/main.rs", r.File)
	assert.Equal(t, uint32(42), r.Line)
	assert.Equal(t, InfoLevel, r.Level)
	assert.Equal(t, "SRC: 1.1.1.1 (22:56:d9:58:18:59), ACTION: DROP", r.Message)

	require.Len(t, r.Args, 5)
	for _, arg := range r.Args {
		assert.NotZero(t, arg.Type)
		assert.NotZero(t, arg.DisplayHint)
		assert.NotEmpty(t, arg.Value)
		assert.NotEmpty(t, arg.Formatted)
	}
}
