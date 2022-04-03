# HDLC

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/borud/hdlc)


HDLC is a library for unframing HDLC-like frames. Instead of implementing this as an io.Reader,
we read whole frames from a channel.  This avoids misunderstandings that might arise if you use the io.Reader interface.

Assuming you have an `io.Reader` named `myReader` this is how you unframe things.

```go
import "github.com/borud/hdlc"
...

 unframer := NewUnframer(myReader)

...

 for range unframer.Frames() {
     // do something with the frames
 }
```

## Framing format

The framing format defines three special values:

```go
const (
    FlagEscape = byte(0x7d)
    FlagSep    = byte(0x7e)
    FlagAbort  = byte(0x7f)
    XORMask    = byte(0x20)
)
```

- Each frame starts and ends with a `FlagSep` byte
- If we encounter a `FlagAbort` value it means we abort the current frame and wait for the next frame beginning before accumulating data.
- If values equal to `FlagEscape`, `FlagSep` or `FlagAbort` occur in the payload it needs to be escaped. Escaping happens by replacing the value with two bytes: FlagEscape followed by the value we are escaping XOR'ed by XORMask (we flip the 5th bit).

Unframing is just applying this process in reverse.