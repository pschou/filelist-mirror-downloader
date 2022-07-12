package main

import (
	"io"

	"github.com/cavaliergopher/rpm"
)

func init() {
	file_parser[".rpm"] = ParseRPMFile
}

func ParseRPMFile(in io.Reader) *file_detail {
	pkg, err := rpm.Read(in)
	if err == nil {
		return &file_detail{
			name: pkg.Name(),
			time: pkg.BuildTime(),
		}
	}
	return nil
}
