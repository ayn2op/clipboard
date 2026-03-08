//go:build ios

package clipboard

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework MobileCoreServices

#import <stdlib.h>
void clipboard_write_string(char *s);
char *clipboard_read_string();
*/
import "C"
import (
	"unsafe"
)

func initialize() error { return nil }

func read(t Format) (buf []byte, err error) {
	switch t {
	case FmtText:
		return []byte(C.GoString(C.clipboard_read_string())), nil
	case FmtImage:
		return nil, ErrUnsupported
	default:
		return nil, ErrUnsupported
	}
}

// write writes the given data to clipboard and
// returns an error if failed.
func write(t Format, buf []byte) error {
	switch t {
	case FmtText:
		cs := C.CString(string(buf))
		defer C.free(unsafe.Pointer(cs))

		C.clipboard_write_string(cs)
		return nil
	case FmtImage:
		return ErrUnsupported
	default:
		return ErrUnsupported
	}
}
