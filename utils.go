package main

import "io"

type nopWriteCloser struct {
	io.Writer
}

func (*nopWriteCloser) Close() {
}
