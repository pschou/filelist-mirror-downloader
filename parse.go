package main

import (
	"io"
	"time"
)

var file_parser = map[string]func(io.ReaderAt) *file_detail{}

type file_detail struct {
	name string
	time time.Time
}
