package main

import (
	"github.com/cavaliergopher/rpm"
)

func init() {
	file_parser[".rpm"] = ParseRPMFile
}

func ParseRPMFile(in ReadAtReader) *file_detail {
	if pkg, err := rpm.Read(in); err == nil {
		return &file_detail{
			name: pkg.Name(),
			time: pkg.BuildTime(),
		}
	}
	return nil
}
