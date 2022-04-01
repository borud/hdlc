# HDLC

HDLC is a library for unframing HDLC-like frames. Instead of implementing this as an io.Reader,
we read whole frames from a channel.  This avoids misunderstandings that might arise if you use the io.Reader interface.

Assuming you have an `io.Reader` named `myReader` this is how you unframe things.

```go
import "github.com/borud/hdlc"
...

 unframer := NewUnframer(bytes.NewReader(myReader))

...

 for range unframer.Frames() {
     // do something with the frames
 }
```
