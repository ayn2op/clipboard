//go:build cgo

package clipboard_test

import (
	"fmt"

	"github.com/ayn2op/clipboard"
)

func ExampleWrite() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	if err := clipboard.Write(clipboard.FmtText, []byte("Hello, 世界")); err != nil {
		panic(err)
	}
	// Output:
}

func ExampleRead() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	data, err := clipboard.Read(clipboard.FmtText)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output:
	// Hello, 世界
}
