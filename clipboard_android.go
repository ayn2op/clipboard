//go:build android

package clipboard

/*
#cgo LDFLAGS: -landroid -llog

#include <stdlib.h>
char *clipboard_read_string(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx);
void clipboard_write_string(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *str);

*/
import "C"
import (
	"unsafe"

	"golang.org/x/mobile/app"
)

func initialize() error { return nil }

func read(t Format) (buf []byte, err error) {
	switch t {
	case FmtText:
		s := ""
		if err := app.RunOnJVM(func(vm, env, ctx uintptr) error {
			cs := C.clipboard_read_string(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx))
			if cs == nil {
				return nil
			}

			s = C.GoString(cs)
			C.free(unsafe.Pointer(cs))
			return nil
		}); err != nil {
			return nil, err
		}
		return []byte(s), nil
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

		if err := app.RunOnJVM(func(vm, env, ctx uintptr) error {
			C.clipboard_write_string(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), cs)
			return nil
		}); err != nil {
			return err
		}
		return nil
	case FmtImage:
		return ErrUnsupported
	default:
		return ErrUnsupported
	}
}
