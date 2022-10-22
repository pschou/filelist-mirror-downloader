package main

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/pschou/go-rpm"
	"github.com/pschou/go-tease"
	"golang.org/x/crypto/openpgp"
)

func init() {
	file_parser[".rpm"] = ParseRPMFile
	file_gpg_check[".rpm"] = CheckRPMFile
}

func ParseRPMFile(in ReadAtReader) *file_detail {
	r := rpm.NewReader(in)

	var err error
	var lead *rpm.Lead
	if lead, err = r.Lead(); err != nil {
		if *debug {
			fmt.Println("Unable to read lead", err)
		}
		return nil
	}

	if lead.SignatureType != 5 {
		if *debug {
			fmt.Println("Unknown signature type:", lead.SignatureType)
		}
		return nil
	}

	var buildTime time.Time

	for i := 0; i < 2; i++ {
		hdr, err := r.Next()
		if err != nil {
			return nil
		}
		for _, t := range hdr.Tags {
			if t.Tag == rpm.RPMTAG_BUILDTIME {
				if iv, ok := t.Int32(); ok && len(iv) > 0 {
					buildTime = time.Unix(int64(iv[0]), 0)
				}
			}
		}
	}

	if !buildTime.IsZero() {
		return &file_detail{
			time: buildTime,
		}
	}
	return nil
}

func CheckRPMFile(fi io.Reader, name string, result chan bool) {
	returnVal := false
	defer func() {
		io.Copy(io.Discard, fi)
		result <- returnVal
	}()
	tr := tease.NewReader(fi)
	r := rpm.NewReader(tr)

	var err error
	var lead *rpm.Lead
	if lead, err = r.Lead(); err != nil {
		if *debug {
			fmt.Println("Error reading rpm lead:", err)
		}
		return
	}

	if lead.SignatureType != 5 {
		if *debug {
			fmt.Println("Unknown RPM signature type:", lead.SignatureType)
		}
		return
	}

	var pgpData []byte

	// Read off the first Header from the RPM
	hdr, err := r.Next()
	if err != nil {
		if *debug {
			fmt.Println("Error parsing RPM header:", err)
		}
		return
	}
	for _, t := range hdr.Tags {
		if t.Count <= 2 {
			continue
		}
		switch t.Tag {
		case rpm.RPMSIGTAG_PGP:
			pgpData, _ = t.Bytes()
			//case rpm.RPMTAG_RSAHEADER:
			//  rsaData, _ = t.Bytes()
		}
	}

	if len(pgpData) == 0 {
		fmt.Println("Signature missing on", name, "- allowing.")
		returnVal = true
		return
	}

	// Align on the 8 byte, and then hash the rest of the file
	offset, err := tr.Seek(0, io.SeekCurrent)
	if err != nil {
		return
	}
	i := (offset + 0x7) &^ 0x7
	_, err = tr.Seek(i, io.SeekStart)
	if err != nil {
		return
	}

	tr.Pipe()

	signer, err := openpgp.CheckDetachedSignature(keyring, fi, bytes.NewReader(pgpData))
	if err != nil {
		return
	}
	returnVal = signer != nil
}
