package hdlc

// http://www.yacer.com.cn/Release/Products/HDLC-UART/HDLC-UART_Datasheet.pdf page 23/24

import (
	"io"
	"log"
)

// Constants for the escape, separator and abort values
const (
	FlagEscape = byte(0x7d)
	FlagSep    = byte(0x7e)
	FlagAbort  = byte(0x7f)
	XORMask    = byte(0x20)
)

const (
	readBufferSize = 512
)

// Unframer unframes HDLC-like frames.
type Unframer struct {
	wrapped io.Reader
	frameCh chan []byte
	err     error
}

const (
	defaultMaxFrameSize = 1024
)

// NewUnframer creates a new unframer.
func NewUnframer(r io.Reader) *Unframer {
	rr := &Unframer{
		wrapped: r,
		err:     nil,
		frameCh: make(chan []byte),
	}

	go rr.readLoop()
	return rr
}

// Frames returns the read channel that returns our frames
func (r *Unframer) Frames() <-chan []byte {
	return r.frameCh
}

func (r *Unframer) Error() error {
	return r.err
}

func (r *Unframer) readLoop() {
	defer func() {
		close(r.frameCh)
		log.Printf("readLoop exited")
	}()

	buf := make([]byte, readBufferSize)
	frame := []byte{}

	// read until we get an EOF or some other error
	for {
		n, err := r.wrapped.Read(buf)
		if err == io.EOF {
			return
		}

		if err != nil {
			r.err = err
			return
		}

		// Undefined what this does, but we abort.
		if n == 0 {
			return
		}

		// scan over the bytes we got
		for _, b := range buf[:n] {
			switch b {
			// Hit a frame separator char
			case FlagSep:
				// If the length of frame is greter than zero this is a frame end.
				if len(frame) > 0 {
					r.frameCh <- Unescape(frame[:])
					frame = []byte{}
				}

			// Abort resets the buffer
			case FlagAbort:
				frame = []byte{}

			default:
				frame = append(frame, b)
			}
		}
	}
}

// Unescape returns an unescaped version of b.
func Unescape(b []byte) []byte {
	unescaped := []byte{}

	xorNext := false
	for _, b := range b {
		// previous byte was FlagEscape so this byte needs to have its fifth bit flipped
		if xorNext {
			unescaped = append(unescaped, b^XORMask)
			xorNext = false
			continue
		}

		// if we hit a FlagEscape we ditch this byte and set xorNext
		if b == FlagEscape {
			xorNext = true
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
			escaped = append(escaped, FlagEscape, FlagEscape^XORMask)

		case FlagAbort:
			escaped = append(escaped, FlagEscape, FlagAbort^XORMask)

		case FlagSep:
			escaped = append(escaped, FlagEscape, FlagSep^XORMask)

		default:
			escaped = append(escaped, b)
		}
	}

	return escaped
}
