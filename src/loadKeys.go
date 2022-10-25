// Written by Paul Schou (paulschou.com) March 2022
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/openpgp"
)

// key rings for signature verfication
var pgpKeys openpgp.EntityList
var rsaKeys []*rsa.PublicKey

func loadKeys(keyfile []byte) (pgpKeys openpgp.EntityList, rsaKeys []*rsa.PublicKey, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(keyfile))
	var line, keystr string
	var i int
	for {
		if scanner.Scan() {
			line = scanner.Text()
		} else {
			break
		}
		keystr += line + "\n"
		switch strings.TrimSpace(line) {
		case "-----END PGP PUBLIC KEY BLOCK-----":
			i++
			var loaded_keys openpgp.EntityList
			loaded_keys, err = openpgp.ReadArmoredKeyRing(strings.NewReader(keystr))
			if err == nil {
				for _, key := range loaded_keys {
					pgpKeys = append(pgpKeys, key)
					fmt.Printf("  %d) Loaded Primary Key (0x%02X)\n", i, key.PrimaryKey.KeyId)
					for _, subkey := range key.Subkeys {
						fmt.Printf("     Sub Key (0x%02X)\n", subkey.PublicKey.KeyId)
					}
				}
				keystr = ""
			} else {
				fmt.Printf("  %d) Invalid key: %g\n", i, err)
			}
		case "-----END PUBLIC KEY-----":
			i++
			var pub_key *pem.Block
			pub_key, _ = pem.Decode([]byte(keystr))
			if pub_key == nil {
				continue
			}

			var key interface{}
			key, err = x509.ParsePKIXPublicKey(pub_key.Bytes)
			if err != nil {
				continue
			}

			if pubKey, ok := key.(*rsa.PublicKey); ok {
				rsaKeys = append(rsaKeys, pubKey)
			}
		}
	}
	if len(pgpKeys) > 0 || len(rsaKeys) > 0 {
		err = nil
	}
	//for _, entity := range []*openpgp.Entity(keyring) {
	//	fmt.Printf("Loaded KeyID: 0x%02X\n", entity.PrimaryKey.KeyId)
	//}
	return
}

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) (exist bool, isdir bool) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, false
	}
	return true, fileInfo.IsDir()
}

func getFiles(walkdir string, suffixes []string) []string {
	ret := []string{}
	err := filepath.Walk(walkdir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}
			for _, suffix := range suffixes {
				if !info.IsDir() && strings.HasSuffix(path, suffix) {
					ret = append(ret, path)
				}
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return ret
}
