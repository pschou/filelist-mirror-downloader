package main

import (
	"archive/tar"
	"compress/gzip"
)

func init() {
	file_parser[".apk"] = ParseAPKFile
}

func ParseAPKFile(in ReadAtReader) *file_detail {
	if gzr, err := gzip.NewReader(in); err == nil {
		defer gzr.Close()
		gzr.Multistream(false)
		if tf := tar.NewReader(gzr); tf != nil {
			if hdr, err := tf.Next(); err == nil {
				return &file_detail{
					time: hdr.ModTime,
				}
			}
		}
	}
	return nil
}
