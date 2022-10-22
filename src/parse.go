package main

import (
	"io"
	"time"
)

var file_parser = map[string]func(ReadAtReader) *file_detail{}

// Return structure for the parsed header
type file_detail struct {
	name string
	time time.Time
}

// General interface which implements both Reader and ReaderAt
type ReadAtReader interface {
	io.Reader
	io.ReaderAt
}

var file_gpg_check = map[string]func(io.Reader, string, chan bool){}
