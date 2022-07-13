package main

import (
	"time"

	"pault.ag/go/debian/deb"
	//"github.com/pschou/go-debian/deb"
)

func init() {
	file_parser[".deb"] = ParseDEBFile
}

func ParseDEBFile(in ReadAtReader) *file_detail {
	if ar, err := deb.LoadAr(in); err == nil {
		if entry, err := ar.Next(); err == nil {
			return &file_detail{
				name: entry.Name,
				time: time.Unix(entry.Timestamp, 0),
			}
		}
	}
	return nil
}
