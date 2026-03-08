/*
Package clipboard provides cross platform clipboard access and supports
macOS/Linux/Windows/Android/iOS platform. Before interacting with the
clipboard, one must call Init to assert if it is possible to use this
package:

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

The most common operations are `Read` and `Write`. To use them:

	// write/read text format data of the clipboard, and
	// the byte buffer regarding the text are UTF8 encoded.
	clipboard.Write(clipboard.FmtText, []byte("text data"))
	clipboard.Read(clipboard.FmtText)

	// write/read image format data of the clipboard, and
	// the byte buffer regarding the image are PNG encoded.
	clipboard.Write(clipboard.FmtImage, []byte("image data"))
	clipboard.Read(clipboard.FmtImage)

Note that read/write regarding image format assumes that the bytes are
PNG encoded since it serves the alpha blending purpose that might be
used in other graphical software.
*/
package clipboard // import "github.com/ayn2op/clipboard"

import (
	"errors"
	"sync"
)

var (
	ErrUnavailable = errors.New("clipboard unavailable")
	ErrUnsupported = errors.New("unsupported format")
	ErrNoCgo       = errors.New("clipboard: cannot use when CGO_ENABLED=0")
)

// Format represents the format of clipboard data.
type Format int

// All sorts of supported clipboard data
const (
	// FmtText indicates plain text clipboard format
	FmtText Format = iota
	// FmtImage indicates image/png clipboard format
	FmtImage
)

var (
	// Due to the limitation on operating systems (such as darwin),
	// concurrent read can even cause panic, use a global lock to
	// guarantee one read at a time.
	lock      = sync.Mutex{}
	initOnce  sync.Once
	initError error
)

// Init initializes the clipboard package. It returns an error
// if the clipboard is not available to use. This may happen if the
// target system lacks required dependency, such as libx11-dev in X11
// environment. For example,
//
//	err := clipboard.Init()
//	if err != nil {
//		panic(err)
//	}
//
// If Init returns an error, any subsequent Read/Write call
// may result in an unrecoverable panic.
func Init() error {
	initOnce.Do(func() {
		initError = initialize()
	})
	return initError
}

// Read returns a chunk of bytes of the clipboard data in the
// desired format. Returns an error if the read fails.
func Read(t Format) ([]byte, error) {
	lock.Lock()
	defer lock.Unlock()
	return read(t)
}

// Write writes a given buffer to the clipboard in a specified format.
// If format t indicates an image, then the given buf assumes
// the image data is PNG encoded.
func Write(t Format, buf []byte) error {
	lock.Lock()
	defer lock.Unlock()
	return write(t, buf)
}
