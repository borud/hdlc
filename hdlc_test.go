package hdlc

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed captured.data
var data []byte

// escapeTests contains escaped and unescaped pairs that we use to perform
// escape and unescape tests.
var escapeTests = []struct {
	Escaped   []byte
	Unescaped []byte
}{
	{
		Escaped:   []byte{FlagEscape, FlagEscape ^ XORMask, 0x99},
		Unescaped: []byte{FlagEscape, 0x99},
	},
	{
		Escaped:   []byte{0x0, FlagEscape, FlagAbort ^ XORMask, 0x99},
		Unescaped: []byte{0x0, FlagAbort, 0x99},
	},
	{
		Escaped:   []byte{FlagEscape, FlagSep ^ XORMask, 0x99},
		Unescaped: []byte{FlagSep, 0x99},
	},
	{
		Escaped:   []byte{FlagEscape, FlagEscape ^ XORMask, FlagEscape, FlagSep ^ XORMask, FlagEscape, FlagAbort ^ XORMask, 0x99},
		Unescaped: []byte{FlagEscape, FlagSep, FlagAbort, 0x99},
	},
}

func TestEscape(t *testing.T) {
	for _, et := range escapeTests {
		assert.Equal(t, et.Unescaped, Unescape(et.Escaped))
		assert.Equal(t, et.Escaped, Escape(et.Unescaped))
		assert.Equal(t, et.Unescaped, Unescape(Escape(et.Unescaped)))
	}
}

func TestHDLC(t *testing.T) {
	unf := NewUnframer(bytes.NewReader(data))

	// for the slice of testdata we're using there should be 9 frames
	numFrames := 0
	for range unf.Frames() {
		numFrames++
	}
	assert.Equal(t, 51, numFrames)
}
