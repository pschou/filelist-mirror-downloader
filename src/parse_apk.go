package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	tease "github.com/pschou/go-tease"
)

func init() {
	file_parser[".apk"] = ParseAPKFile
	file_sig_check[".apk"] = CheckAPKFile
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

func CheckAPKFile(fi io.Reader, name string, result chan bool) {
	returnVal := false
	defer func() {
		io.Copy(io.Discard, fi)
		result <- returnVal
	}()
	tr := tease.NewReader(fi)
	gzr, err := gzip.NewReader(tr)
	if err != nil {
		if *debug {
			fmt.Errorf("Invalid APK file,", err)
		}
		return
	}

	defer gzr.Close()
	gzr.Multistream(false)

	tarRdr := tar.NewReader(gzr)
	header, err := tarRdr.Next()
	if err != nil {
		if *debug {
			fmt.Errorf("Invalid APK tar header,", err)
		}
		return
	}

	var startOffset int64
	var sig []byte
	if strings.HasPrefix(header.Name, ".SIGN.RSA.") {
		sig, err = ioutil.ReadAll(tarRdr)
		if len(sig) == 0 {
			return
		}
		// Read to the end of the first gzip section
		_, err = io.Copy(io.Discard, gzr)
		if err != nil {
			if *debug {
				fmt.Errorf("Invalid APK tar body,", err)
			}
			return
		}

		// Get the current pointer in the file
		startOffset, _ = tr.Seek(0, io.SeekCurrent)

		// Start the new section
		gzr.Reset(tr)
		gzr.Multistream(false)
	}

	// Read to the end of the second gzip section
	_, err = io.Copy(io.Discard, gzr)
	if err != nil {
		if *debug {
			fmt.Errorf("Invalid APK second tar body,", err)
		}
		return
	}

	// Get the current pointer in the file
	endOffset, _ := tr.Seek(0, io.SeekCurrent)

	if endOffset <= startOffset {
		return
	}

	// Read in second section to check signature
	tr.Seek(startOffset, io.SeekStart)
	tr.Pipe()
	hash := sha1.New()
	io.Copy(hash, io.LimitReader(tr, endOffset-startOffset))

	for i, pubKey := range rsaKeys {
		err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA1, hash.Sum(nil)[:], sig)
		if err == nil {
			if *debug {
				fmt.Errorf("Valid APK file,", i)
			}
			returnVal = true
		}
	}
	if *debug {
		fmt.Errorf("Invalid APK file")
	}
}
