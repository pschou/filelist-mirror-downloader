package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	tease "github.com/pschou/go-tease"
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

	//buff            *bytes.Buffer
	reader *io.PipeReader
	writer *io.PipeWriter
	err    error

	hashPKGINFO     hash.Hash
	expectedPKGINFO string
	dataHash        hash.Hash
	expectedDATA    string
	//inData          *bool
	mutex sync.Mutex
}

func NewWithExpectedHash(hash string) hash.Hash {
	//fmt.Println("built out with expected hash", hash)
	r, w := io.Pipe()

	a := alpineHash{
		reader: r,
		writer: w,
		//buff:            bytes.NewBuffer([]byte{}),
		hashPKGINFO:     sha1.New(),
		dataHash:        sha256.New(),
		expectedPKGINFO: hash,
	}
	go a.parse()
	return a
}

func (a alpineHash) Write(p []byte) (n int, err error) {
	if a.err != nil {
		return 0, a.err
	}
	return a.writer.Write(p)
}

func (a *alpineHash) parse() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	tr := tease.NewReader(a.reader)
	gzr, err := gzip.NewReader(tr)
	if err != nil {
		a.err = err
		return
	}
	defer gzr.Close()
	gzr.Multistream(false)

	tarRdr := tar.NewReader(gzr)
	header, err := tarRdr.Next()
	if err != nil {
		a.err = err
		return
	}

	var startOffset int64
	if strings.HasPrefix(header.Name, ".SIGN.RSA.") {
		//fmt.Println("found sign rsa section")
		// Read to the end of the first gzip section
		io.Copy(io.Discard, gzr)

		// Get the current pointer in the file
		startOffset, _ = tr.Seek(0, io.SeekCurrent)

		// Start the new section
		gzr.Reset(tr)
		gzr.Multistream(false)

		//fmt.Println("rebuild gzip")
		tarRdr = tar.NewReader(gzr)
		header, err = tarRdr.Next()
		if err != nil {
			a.err = err
			//fmt.Println("advancing tar", err)
			return
		}
	}

	//fmt.Println("found", header)
	if strings.HasPrefix(header.Name, ".PKGINFO") {
		//fmt.Println("found pkginfo")
		txt, err := ioutil.ReadAll(tarRdr)
		if len(txt) == 0 || err != nil {
			a.err = fmt.Errorf("invalid pkginfo")
			//fmt.Println("invalid pkginfo", err)
			return
		}

		// Extract the payload datahash
		pos_hash := strings.Index(string(txt), "datahash = ")
		if pos_hash >= 0 && len(txt) > pos_hash+75 {
			a.expectedDATA = string(txt[pos_hash+11 : pos_hash+75])
			//fmt.Println("datahash =", a.expectedDATA)
		} else {
			a.err = fmt.Errorf("missing datahash")
			fmt.Println("missing datahash")
			return
		}

		// Read to the end of the second gzip section
		io.Copy(io.Discard, gzr)

		// Get the current pointer in the file
		endOffset, _ := tr.Seek(0, io.SeekCurrent)

		// Read in the raw bytes to our hash
		tr.Seek(startOffset, io.SeekStart)
		io.Copy(a.hashPKGINFO, io.LimitReader(tr, endOffset-startOffset))

		// Compare the PKGINFO with the expected Q1 value
		s := a.hashPKGINFO.Sum(nil)
		// For some odd reason, sometimes this is all lower cased:
		if !strings.EqualFold(a.expectedPKGINFO, "Q1"+base64.StdEncoding.EncodeToString(s)) {
			if *debug {
				//fmt.Printf("  %v != %v\n", a.expectedPKGINFO, "Q1"+base64.StdEncoding.EncodeToString(s))
			}
			a.err = fmt.Errorf("PKGINFO hash mismatch in APK file")
			return
		}
	}

	tr.Pipe()
	io.Copy(a.dataHash, tr)
}

/*func gunzip(dat []byte) (string, error) {
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
}*/

func (a alpineHash) Sum(b []byte) []byte {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	//a.Write(b)
	s := a.hashPKGINFO.Sum(nil)
	if a.expectedDATA != "" {
		if !strings.EqualFold(a.expectedDATA, fmt.Sprintf("%x", a.dataHash.Sum(nil))) {
			return []byte{}
		}
	}
	return []byte("Q1" + base64.StdEncoding.EncodeToString(s))
}

func (a alpineHash) Reset() {}

func (a alpineHash) Size() int {
	return 30
}
func (a alpineHash) BlockSize() int {
	return 512
}

//var gzip_header = []byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00}

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
