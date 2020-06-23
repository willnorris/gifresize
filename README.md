# gifresize 

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/willnorris.com/go/gifresize)

gifresize is a simple go package for transforming animated GIFs.

Import using:

```go
import "willnorris.com/go/gifresize"
```

Then call `gifresize.Process` with the source io.Reader and destination
io.Writer as well as the transformation to be applied to each frame in the GIF.
See [example/main.go][] for a simple example.

[example/main.go]: ./example/main.go
