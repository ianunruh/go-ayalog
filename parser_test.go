package ayalog

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRecord_Simple(t *testing.T) {
	record, err := os.ReadFile("testdata/simple.bin")
	require.NoError(t, err)

	parser := Parser{
		IncludeArgs: true,
	}

	r, err := parser.Record(bytes.NewBuffer(record))
	require.NoError(t, err)

	assert.Equal(t, "xdp_hello", r.Target)
	assert.Equal(t, "xdp_hello", r.Module)
	assert.Equal(t, "src/main.rs", r.File)
	assert.Equal(t, uint32(40), r.Line)
	assert.Equal(t, InfoLevel, r.Level)
	assert.Equal(t, "SRC: 1.1.1.1, ACTION: DROP", r.Message)

	require.Len(t, r.Args, 3)
	for _, arg := range r.Args {
		assert.NotZero(t, arg.Type)
		assert.NotZero(t, arg.DisplayHint)
		assert.NotEmpty(t, arg.Value)
		assert.NotEmpty(t, arg.Formatted)
	}
}
