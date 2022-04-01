package hdlc

import (
	"io"
)

// Constants for the escape, separator and abort values
const (
	FlagEscape = byte(0x7d)
	FlagSep    = byte(0x7e)
	FlagAbort  = byte(0x7f)
)

// Unframer unframes HDLC-like frames.
type Unframer struct {
	wrapped         io.Reader
	frameCh         chan []byte
	skipEmptyFrames bool
	err             error
	maxFrameSize    int
}

const (
	defaultMaxFrameSize = 1024
)

// NewUnframer creates a new unframer.
func NewUnframer(r io.Reader) *Unframer {
	rr := &Unframer{
		wrapped:         r,
		skipEmptyFrames: true,
		err:             nil,
		maxFrameSize:    defaultMaxFrameSize,
		frameCh:         make(chan []byte),
	}

	go rr.readLoop()

	return rr
}

// SetSkipEmptyFrames configures if you skip empty frames or return them. If set to true,
// the Unframer will skip empty frames.  (The default is that we skip empty frames)
func (r *Unframer) SetSkipEmptyFrames(skip bool) *Unframer {
	r.skipEmptyFrames = skip
	return r
}

// Frames returns the read channel that returns our frames
func (r *Unframer) Frames() <-chan []byte {
	return r.frameCh
}

func (r *Unframer) Error() error {
	return r.err
}

func (r *Unframer) readLoop() {
	defer close(r.frameCh)

	buf := make([]byte, r.maxFrameSize+1)
	frame := []byte{}

	// read until we get an EOF or some other error
	for {
		_, err := r.wrapped.Read(buf)
		if err == io.EOF {
			return
		}

		if err != nil {
			r.err = err
			return
		}

		// scan over the bytes we got
		for _, b := range buf {
			if b == FlagSep {
				// Skip empty frames
				if r.skipEmptyFrames && len(frame) == 0 {
					continue
				}

				// unescape and return frame
				r.frameCh <- Unescape(frame)
				frame = []byte{}
				continue
			}

			if b == FlagAbort {
				frame = []byte{}
				continue
			}

			frame = append(frame, b)

		}
	}
}

// Unescape returns an unescaped version of b.
func Unescape(b []byte) []byte {
	unescaped := []byte{}

	flipBitFiveOnNextByte := false
	for _, b := range b {
		// previous byte was FlagEscape so this byte needs to have its fifth bit flipped
		if flipBitFiveOnNextByte {
			unescaped = append(unescaped, b^0x20)
			flipBitFiveOnNextByte = false
			continue
		}

		// if we hit a FlagEscape we ditch this byte and set flipBitFiveOnNextByte
		if b == FlagEscape {
			flipBitFiveOnNextByte = true
			continue
		}

		// otherwise we just append
		unescaped = append(unescaped, b)
	}

	return unescaped
}

// Escape returns an escaped version of b.
func Escape(b []byte) []byte {
	escaped := []byte{}
	for _, b := range b {
		switch b {
		case FlagEscape:
			escaped = append(escaped, FlagEscape, FlagEscape^0x20)

		case FlagAbort:
			escaped = append(escaped, FlagEscape, FlagAbort^0x20)

		case FlagSep:
			escaped = append(escaped, FlagEscape, FlagSep^0x20)

		default:
			escaped = append(escaped, b)
		}
	}

	return escaped
}
