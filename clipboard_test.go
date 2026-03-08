package clipboard_test

import (
	"bytes"
	"errors"
	"image/color"
	"image/png"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/ayn2op/clipboard"
)

func init() {
	// No debug mode needed anymore
}

func TestClipboardInit(t *testing.T) {
	t.Run("no-cgo", func(t *testing.T) {
		if val, ok := os.LookupEnv("CGO_ENABLED"); !ok || val != "0" {
			t.Skip("CGO_ENABLED is set to 1")
		}
		if runtime.GOOS == "windows" {
			t.Skip("Windows does not need to check for cgo")
		}

		if err := clipboard.Init(); !errors.Is(err, clipboard.ErrNoCgo) {
			t.Fatalf("expect ErrNoCgo, got: %v", err)
		}
	})
	t.Run("with-cgo", func(t *testing.T) {
		if val, ok := os.LookupEnv("CGO_ENABLED"); ok && val == "0" {
			t.Skip("CGO_ENABLED is set to 0")
		}
		if runtime.GOOS != "linux" && runtime.GOOS != "freebsd" {
			t.Skip("Only Linux/FreeBSD may return error at the moment.")
		}

		if err := clipboard.Init(); err != nil && !errors.Is(err, clipboard.ErrUnavailable) {
			t.Fatalf("expect ErrUnavailable, but got: %v", err)
		}
	})
}

func TestClipboard(t *testing.T) {
	if runtime.GOOS != "windows" {
		if val, ok := os.LookupEnv("CGO_ENABLED"); ok && val == "0" {
			t.Skip("CGO_ENABLED is set to 0")
		}
	}

	t.Run("image", func(t *testing.T) {
		data, err := os.ReadFile("tests/testdata/clipboard.png")
		if err != nil {
			t.Fatalf("failed to read gold file: %v", err)
		}
		if err := clipboard.Write(clipboard.FmtImage, data); err != nil {
			t.Fatalf("failed to write image: %v", err)
		}

		_, err = clipboard.Read(clipboard.FmtText)
		if err == nil {
			t.Fatalf("read clipboard that stores image data as text should fail, but succeeded")
		}

		b, err := clipboard.Read(clipboard.FmtImage)
		if err != nil {
			t.Fatalf("failed to read image: %v", err)
		}
		if b == nil {
			t.Fatalf("read clipboard that stores image data as image should succeed, but got: nil")
		}

		img1, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("write image is not png encoded: %v", err)
		}
		img2, err := png.Decode(bytes.NewReader(b))
		if err != nil {
			t.Fatalf("read image is not png encoded: %v", err)
		}

		w := img2.Bounds().Dx()
		h := img2.Bounds().Dy()

		incorrect := 0
		for i := range w {
			for j := range h {
				wr, wg, wb, wa := img1.At(i, j).RGBA()
				gr, gg, gb, ga := img2.At(i, j).RGBA()
				want := color.RGBA{
					R: uint8(wr),
					G: uint8(wg),
					B: uint8(wb),
					A: uint8(wa),
				}
				got := color.RGBA{
					R: uint8(gr),
					G: uint8(gg),
					B: uint8(gb),
					A: uint8(ga),
				}

				if !reflect.DeepEqual(want, got) {
					t.Logf("read data from clipboard is inconsistent with previous written data, pix: (%d,%d), got: %+v, want: %+v", i, j, got, want)
					incorrect++
				}
			}
		}

		if incorrect > 0 {
			t.Fatalf("read data from clipboard contains too much inconsistent pixels to the previous written data, number of incorrect pixels: %v", incorrect)
		}
	})

	t.Run("text", func(t *testing.T) {
		data := []byte("github.com/ayn2op/clipboard")
		if err := clipboard.Write(clipboard.FmtText, data); err != nil {
			t.Fatalf("failed to write text: %v", err)
		}

		_, err := clipboard.Read(clipboard.FmtImage)
		if err == nil {
			t.Fatalf("read clipboard that stores text data as image should fail, but succeeded")
		}

		b, err := clipboard.Read(clipboard.FmtText)
		if err != nil {
			t.Fatalf("failed to read text: %v", err)
		}
		if b == nil {
			t.Fatal("read clipboard that stores text data as text should succeed, but got: nil")
		}

		if !reflect.DeepEqual(data, b) {
			t.Fatalf("read data from clipbaord is inconsistent with previous written data, got: %d, want: %d", len(b), len(data))
		}
	})
}

func TestClipboardMultipleWrites(t *testing.T) {
	if runtime.GOOS != "windows" {
		if val, ok := os.LookupEnv("CGO_ENABLED"); ok && val == "0" {
			t.Skip("CGO_ENABLED is set to 0")
		}
	}

	data, err := os.ReadFile("tests/testdata/clipboard.png")
	if err != nil {
		t.Fatalf("failed to read gold file: %v", err)
	}
	if err := clipboard.Write(clipboard.FmtImage, data); err != nil {
		t.Fatalf("failed to write image: %v", err)
	}

	data = []byte("github.com/ayn2op/clipboard")
	if err := clipboard.Write(clipboard.FmtText, data); err != nil {
		t.Fatalf("failed to write text: %v", err)
	}

	b, err := clipboard.Read(clipboard.FmtImage)
	if err == nil && b != nil {
		t.Fatalf("read clipboard that should store text data as image should fail, but got: %d", len(b))
	}

	b, err = clipboard.Read(clipboard.FmtText)
	if err != nil {
		t.Fatalf("failed to read text: %v", err)
	}
	if b == nil {
		t.Fatalf("read clipboard that should store text data as text should succeed, got: nil")
	}

	if !reflect.DeepEqual(data, b) {
		t.Fatalf("read data from clipboard is inconsistent with previous write, want %s, got: %s", string(data), string(b))
	}
}

func TestClipboardConcurrentRead(t *testing.T) {
	if runtime.GOOS != "windows" {
		if val, ok := os.LookupEnv("CGO_ENABLED"); ok && val == "0" {
			t.Skip("CGO_ENABLED is set to 0")
		}
	}

	// This test check that concurrent read/write to the clipboard does
	// not cause crashes on some specific platform, such as macOS.
	done := make(chan bool, 2)
	go func() {
		defer func() {
			done <- true
		}()
		clipboard.Read(clipboard.FmtText)
	}()
	go func() {
		defer func() {
			done <- true
		}()
		clipboard.Read(clipboard.FmtImage)
	}()
	<-done
	<-done
}

func TestClipboardWriteEmpty(t *testing.T) {
	if runtime.GOOS != "windows" {
		if val, ok := os.LookupEnv("CGO_ENABLED"); ok && val == "0" {
			t.Skip("CGO_ENABLED is set to 0")
		}
	}

	if err := clipboard.Write(clipboard.FmtText, nil); err != nil {
		t.Fatalf("failed to write nil: %v", err)
	}
	got, err := clipboard.Read(clipboard.FmtText)
	if err != nil {
		t.Fatalf("failed to read text: %v", err)
	}
	if got != nil {
		t.Fatalf("write nil to clipboard should read nil, got: %v", string(got))
	}

	if err := clipboard.Write(clipboard.FmtText, []byte("")); err != nil {
		t.Fatalf("failed to write empty string: %v", err)
	}
	got, err = clipboard.Read(clipboard.FmtText)
	if err != nil {
		t.Fatalf("failed to read text: %v", err)
	}
	if string(got) != "" {
		t.Fatalf("write empty string to clipboard should read empty string, got: `%v`", string(got))
	}
}

func BenchmarkClipboard(b *testing.B) {
	b.Run("text", func(b *testing.B) {
		data := []byte("github.com/ayn2op/clipboard")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := clipboard.Write(clipboard.FmtText, data); err != nil {
				b.Fatalf("failed to write: %v", err)
			}
			if _, err := clipboard.Read(clipboard.FmtText); err != nil {
				b.Fatalf("failed to read: %v", err)
			}
		}
	})
}

func TestClipboardNoCgo(t *testing.T) {
	if val, ok := os.LookupEnv("CGO_ENABLED"); !ok || val != "0" {
		t.Skip("CGO_ENABLED is set to 1")
	}
	if runtime.GOOS == "windows" {
		t.Skip("Windows should always be tested")
	}

	t.Run("Read", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				return
			}
			t.Fatalf("expect to fail when CGO_ENABLED=0")
		}()

		clipboard.Read(clipboard.FmtText)
	})

	t.Run("Write", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				return
			}
			t.Fatalf("expect to fail when CGO_ENABLED=0")
		}()

		clipboard.Write(clipboard.FmtText, []byte("dummy"))
	})
}
