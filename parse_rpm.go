package main

import (
	"io"

	"github.com/cavaliergopher/rpm"
)

func init() {
	file_parser[".rpm"] = ParseRPMFile
}

func ParseRPMFile(in io.ReaderAt) *file_detail {
	pkg, err := rpm.Read(in.(io.Reader))
	if err == nil {
		return &file_detail{
			name: pkg.Name(),
			time: pkg.BuildTime(),
		}
	}
	return nil
}
