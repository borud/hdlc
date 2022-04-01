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
		Escaped:   []byte{FlagEscape, FlagEscape ^ 0x20, 0x99},
		Unescaped: []byte{FlagEscape, 0x99},
	},
	{
		Escaped:   []byte{0x0, FlagEscape, FlagAbort ^ 0x20, 0x99},
		Unescaped: []byte{0x0, FlagAbort, 0x99},
	},
	{
		Escaped:   []byte{FlagEscape, FlagSep ^ 0x20, 0x99},
		Unescaped: []byte{FlagSep, 0x99},
	},
	{
		Escaped:   []byte{FlagEscape, FlagEscape ^ 0x20, FlagEscape, FlagSep ^ 0x20, FlagEscape, FlagAbort ^ 0x20, 0x99},
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
	unf := NewUnframer(bytes.NewReader(data[:200]))

	// for the slice of testdata we're using there should be 9 frames
	numFrames := 0
	for range unf.Frames() {
		numFrames++
	}
	assert.Equal(t, 9, numFrames)
}

func TestHDLCEmptyFrames(t *testing.T) {
	unf := NewUnframer(bytes.NewReader(data[:200])).
		SetSkipEmptyFrames(false)

	numFrames := 0
	for range unf.Frames() {
		numFrames++
	}
	assert.Equal(t, 19, numFrames)
}
