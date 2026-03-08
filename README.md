# clipboard [![PkgGoDev](https://pkg.go.dev/badge/github.com/ayn2op/clipboard)](https://pkg.go.dev/github.com/ayn2op/clipboard) ![](https://changkun.de/urlstat?mode=github&repo=golang-design/clipboard) ![clipboard](https://github.com/golang-design/clipboard/workflows/clipboard/badge.svg?branch=main)

Cross platform (macOS/Linux/FreeBSD/Windows/Android/iOS) clipboard package in Go

```go
import "github.com/ayn2op/clipboard"
```

## Features

- Cross platform supports: **macOS, Linux (X11), FreeBSD (X11), Windows, iOS, and Android**
- Copy/paste UTF-8 text
- Copy/paste PNG encoded images (Desktop-only)

## API Usage

Package clipboard provides cross platform clipboard access and supports
macOS/Linux/FreeBSD/Windows/Android/iOS platform. Before interacting with the
clipboard, one must call Init to assert if it is possible to use this
package:

```go
// Init returns an error if the package is not ready for use.
err := clipboard.Init()
if err != nil {
      panic(err)
}
```

The most common operations are `Read` and `Write`. To use them:

```go
// write/read text format data of the clipboard, and
// the byte buffer regarding the text are UTF8 encoded.
if err := clipboard.Write(clipboard.FmtText, []byte("text data")); err != nil {
    panic(err)
}
clipboard.Read(clipboard.FmtText)

// write/read image format data of the clipboard, and
// the byte buffer regarding the image are PNG encoded.
if err := clipboard.Write(clipboard.FmtImage, []byte("image data")); err != nil {
    panic(err)
}
clipboard.Read(clipboard.FmtImage)
```

Note that read/write regarding image format assumes that the bytes are
PNG encoded since it serves the alpha blending purpose that might be
used in other graphical software.

## Platform Specific Details

This package spent efforts to provide cross platform abstraction regarding
accessing system clipboards, but here are a few details you might need to know.

### Dependency

- macOS: require Cgo, no dependency
 - Linux/FreeBSD: require X11 dev package. For instance, install `libx11-dev` or `xorg-dev` or `libX11-devel` on Linux, or `libX11` on FreeBSD to access X window system.
   Wayland sessions are currently unsupported; running under Wayland
   typically requires an XWayland bridge and `DISPLAY` to be set.
- Windows: no Cgo, no dependency
- iOS/Android: collaborate with [`gomobile`](https://golang.org/x/mobile)

### Screenshot

In general, when you need test your implementation regarding images,
There are system level shortcuts to put screenshot image into your system clipboard:

- On macOS, use `Ctrl+Shift+Cmd+4`
- On Linux/Ubuntu, use `Ctrl+Shift+PrintScreen`
- On Windows, use `Shift+Win+s`

As described in the API documentation, the package supports read/write
UTF8 encoded plain text or PNG encoded image data. Thus,
the other types of data are not supported yet, i.e. undefined behavior.

## License

MIT | &copy; 2021 The golang.design Initiative Authors, written by [Changkun Ou](https://changkun.de).
