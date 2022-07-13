package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"strings"
)

type alpineHash struct {
	// Write (via the embedded io.Writer interface) adds more data to the running hash.
	// It never returns an error.
	//io.Writer

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	//Sum(b []byte) []byte

	// Reset resets the Hash to its initial state.
	//Reset()

	// Size returns the number of bytes Sum will return.
	//Size() int

	// BlockSize returns the hash's underlying block size.
	// The Write method must be able to accept any amount
	// of data, but it may operate more efficiently if all writes
	// are a multiple of the block size.
	//BlockSize() int

	buff            *bytes.Buffer
	hashPKGINFO     hash.Hash
	expectedPKGINFO *string
	hashDATA        hash.Hash
	expectedDATA    *string
	inData          *bool
}

func NewWithExpectedHash(hash string) hash.Hash {
	return alpineHash{
		buff:            bytes.NewBuffer([]byte{}),
		expectedPKGINFO: &hash,
		hashPKGINFO:     sha1.New(),
		expectedDATA:    new(string),
		hashDATA:        sha256.New(),
		inData:          new(bool),
	}
}

func (a alpineHash) Write(p []byte) (n int, err error) {
	if *a.inData {
		n, err = a.hashDATA.Write(p)
	} else {
		n, err = a.buff.Write(p)

		dat := a.buff.Bytes()
		if !bytes.HasPrefix(dat, gzip_header) {
			return 0, fmt.Errorf("Invalid APK file")
		}

		pos := bytes.Index(dat[1:], gzip_header) + 1
		for pos > 0 && !*a.inData {
			txt, err := gunzip(a.buff.Next(pos))
			if err != nil {
				return 0, err
			}
			switch {
			case strings.HasPrefix(txt, ".SIGN.RSA."):
				//fmt.Println("found rsa")
				dat = a.buff.Bytes()
				if !bytes.HasPrefix(dat, gzip_header) {
					return 0, fmt.Errorf("Invalid APK file")
				}

				pos = bytes.Index(dat[1:], gzip_header) + 1

			case strings.HasPrefix(txt, "./PaxHeaders/.PKGINFO"):
				//fmt.Println("found pkginfo")
				*a.inData = true
				pos_hash := strings.Index(txt, "datahash = ")
				*a.expectedDATA = txt[pos_hash+11 : pos_hash+75]
				a.hashPKGINFO.Reset()
				a.hashPKGINFO.Write(dat[:pos])
				s := a.hashPKGINFO.Sum(nil)
				// For some odd reason, sometimes this is all lower cased:
				if !strings.EqualFold(*a.expectedPKGINFO, "Q1"+base64.StdEncoding.EncodeToString(s)) {
					if *debug {
						fmt.Printf("  %v != %v\n", *a.expectedPKGINFO, "Q1"+base64.StdEncoding.EncodeToString(s))
					}
					return 0, fmt.Errorf("PKGINFO hash mismatch in APK file")
				}
				a.hashDATA.Reset()
				a.hashDATA.Write(a.buff.Next(a.buff.Len()))
				return n, err
			}
		}

		//if a.buff.Len() > 10000 {
		//	return 0, fmt.Errorf("APK file in wrong format, missing header")
		//}
	}
	return
}

func gunzip(dat []byte) (string, error) {
	gz_dat, err := gzip.NewReader(bytes.NewReader(dat))
	if err != nil {
		return "", err
	}
	gz_dat.Close()
	buf := new(strings.Builder)
	_, err = io.Copy(buf, gz_dat)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (a alpineHash) Sum(b []byte) []byte {
	a.Write(b)
	s := a.hashPKGINFO.Sum(nil)
	if a.expectedDATA != nil {
		if !strings.EqualFold(*a.expectedDATA, fmt.Sprintf("%x", a.hashDATA.Sum(nil))) {
			return []byte{}
		}
	}
	return []byte("Q1" + base64.StdEncoding.EncodeToString(s))
}

func (a alpineHash) Reset() {
	a.buff.Reset()
	a.hashPKGINFO.Reset()
	*a.expectedPKGINFO = ""
	a.hashDATA.Reset()
	*a.expectedDATA = ""
	*a.inData = false
}

func (a alpineHash) Size() int {
	return 30
}
func (a alpineHash) BlockSize() int {
	return 512
}

var gzip_header = []byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00}

/*
func split_on_gzip_header(data []byte) ([]byte, []byte) {
	gzip_header := []byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00}

	var pos int
	for pos := 0; pos < len(data)-8; pos++ {
		if bytes.Equal(data[pos:pos+len(gzip_header)], gzip_header) {
			sig := data[:pos]
			content := data[pos:]

			return sig, content
		}
	}

	return []byte{}, []byte{}
}*/
