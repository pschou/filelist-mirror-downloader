package main

import (
	"io"
	"time"

	"pault.ag/go/debian/deb"
	//"github.com/pschou/go-debian/deb"
)

func init() {
	file_parser[".deb"] = ParseDEBFile
}

func ParseDEBFile(in io.ReaderAt) *file_detail {
	ar, err := deb.LoadAr(in)
	if err == nil {
		entry, err := ar.Next()
		if err == nil {
			return &file_detail{
				name: entry.Name,
				time: time.Unix(entry.Timestamp, 0),
			}
		}
	}
	return nil
}
