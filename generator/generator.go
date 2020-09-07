package main

import (
	"fmt"
	"os"
	"strings"
)

func genFunc(fh *os.File, name, typ string) {
	if name == "" {
		if strings.HasPrefix(typ, "[]") {
			name = strings.Title(typ[2:]) + "s"
		} else {
			name = strings.Title(typ)
		}
	}
	fmt.Fprintf(fh, `
// %s panics on error, returns the first argument otherwise
func %s(arg %s, err error) %s {
	OK(err)
	return arg
}
`, name, name, typ, typ)
}

func genFuncs(fh *os.File, typ string) {
	genFunc(fh, "", typ)
	genFunc(fh, "", "[]"+typ)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: generator <output-file>\n")
		os.Exit(2)
	}
	fh, err := os.OpenFile(os.Args[1], os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open file for writing: %v\n", err)
		os.Exit(1)
	}
	defer fh.Close()

	fmt.Fprint(fh, `// Code generated by github.com/ridge/must/generator. DO NOT EDIT.

package must

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
)
`)

	for _, size := range []string{"8", "16", "32", "64", ""} {
		genFuncs(fh, "uint"+size)
		genFuncs(fh, "int"+size)
	}
	for _, typ := range []string{"bool", "float32", "float64",
		"complex64", "complex128", "byte", "rune", "uintptr",
		"string"} {
		genFuncs(fh, typ)
	}
	genFunc(fh, "Any", "interface{}")
	genFunc(fh, "OSFile", "*os.File")
	genFunc(fh, "OSFileInfo", "os.FileInfo")
	genFunc(fh, "OSFileInfos", "[]os.FileInfo")
	genFunc(fh, "IOReadCloser", "io.ReadCloser")
	genFunc(fh, "IOWriter", "io.Writer")
	genFunc(fh, "NetIP", "net.IP")
	genFunc(fh, "NetListener", "net.Listener")
	genFunc(fh, "NetURL", "*url.URL")
	genFunc(fh, "HTTPRequest", "*http.Request")
	genFunc(fh, "HTTPHandler", "http.Handler")
	genFunc(fh, "Time", "time.Time")
}
