//go:build darwin && !ios

package clipboard

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

unsigned int clipboard_read_string(void **out);
unsigned int clipboard_read_image(void **out);
int clipboard_write_string(const void *bytes, NSInteger n);
int clipboard_write_image(const void *bytes, NSInteger n);
*/
import "C"
import (
	"unsafe"
)

func initialize() error { return nil }

func read(t Format) (buf []byte, err error) {
	var (
		data unsafe.Pointer
		n    C.uint
	)
	switch t {
	case FmtText:
		n = C.clipboard_read_string(&data)
	case FmtImage:
		n = C.clipboard_read_image(&data)
	}
	if data == nil {
		return nil, ErrUnavailable
	}
	defer C.free(unsafe.Pointer(data))
	if n == 0 {
		return nil, nil
	}
	return C.GoBytes(data, C.int(n)), nil
}

// write writes the given data to clipboard and returns an error if failed.
func write(t Format, buf []byte) error {
	var ok C.int
	switch t {
	case FmtText:
		if len(buf) == 0 {
			ok = C.clipboard_write_string(unsafe.Pointer(nil), 0)
		} else {
			ok = C.clipboard_write_string(unsafe.Pointer(&buf[0]),
				C.NSInteger(len(buf)))
		}
	case FmtImage:
		if len(buf) == 0 {
			ok = C.clipboard_write_image(unsafe.Pointer(nil), 0)
		} else {
			ok = C.clipboard_write_image(unsafe.Pointer(&buf[0]),
				C.NSInteger(len(buf)))
		}
	default:
		return ErrUnsupported
	}
	if ok != 0 {
		return ErrUnavailable
	}
	return nil
}
