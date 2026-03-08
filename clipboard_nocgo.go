//go:build !windows && !cgo

package clipboard

func initialize() error {
	return ErrNoCgo
}

func read(t Format) (buf []byte, err error) {
	panic("clipboard: cannot use when CGO_ENABLED=0")
}

func readc(t string) ([]byte, error) {
	panic("clipboard: cannot use when CGO_ENABLED=0")
}

func write(t Format, buf []byte) error {
	panic("clipboard: cannot use when CGO_ENABLED=0")
}
