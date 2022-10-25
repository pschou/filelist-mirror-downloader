package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	tease "github.com/pschou/go-tease"
	"golang.org/x/crypto/openpgp"
	"pault.ag/go/debian/deb"
	//"github.com/pschou/go-debian/deb"
)

func init() {
	file_parser[".deb"] = ParseDEBFile
	file_sig_check[".deb"] = CheckDEBFile
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

func CheckDEBFile(fi io.Reader, name string, result chan bool) {
	returnVal := false
	defer func() {
		io.Copy(io.Discard, fi)
		result <- returnVal
	}()
	tr := tease.NewReader(fi)

	if ar, err := deb.LoadAr(tr); err == nil {
		var sig, binaryFlag, control, data *deb.ArEntry

		for entry, err := ar.Next(); err == nil; entry, err = ar.Next() {
			switch {
			case entry.Name == `debian-binary`:
				binaryFlag = entry
			case strings.HasPrefix(entry.Name, "_gpg"):
				fmt.Println("my  entry:", entry.Name)
				sig = entry
			case strings.HasPrefix(entry.Name, "control."):
				control = entry
			case strings.HasPrefix(entry.Name, "data."):
				data = entry
			}
		}
		// Many packages are not signed so doing this may cause nothing to be downloaded
		if sig == nil {
			returnVal = true
			return
		}
		if control == nil || data == nil || binaryFlag == nil {
			// "unable to find signed data"
			return
		}
		binaryFlag.Data.Seek(0, 0)
		control.Data.Seek(0, 0)
		data.Data.Seek(0, 0)

		var r io.Reader
		r = debrdr{tr: tr, r: data.Data}
		signedData := io.MultiReader(binaryFlag.Data, control.Data, r)
		signer, _ := openpgp.CheckDetachedSignature(pgpKeys, signedData, sig.Data)
		returnVal = signer != nil
	}
}

type debrdr struct {
	tr *tease.Reader
	r  io.Reader
}

func (r debrdr) Read(p []byte) (n int, err error) {
	r.tr.Pipe()
	return r.r.Read(p)
}
