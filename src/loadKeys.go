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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/openpgp"
)

func loadKeys(keyfile []byte) (keyring openpgp.EntityList, err error) {
	var loaded_keys openpgp.EntityList
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
		if strings.TrimSpace(line) == "-----END PGP PUBLIC KEY BLOCK-----" {
			i++
			loaded_keys, err = openpgp.ReadArmoredKeyRing(strings.NewReader(keystr))
			if err == nil {
				for _, key := range loaded_keys {
					keyring = append(keyring, key)
					fmt.Printf("  %d) Loaded Primary Key (0x%02X)\n", i, key.PrimaryKey.KeyId)
					for _, subkey := range key.Subkeys {
						fmt.Printf("     Sub Key (0x%02X)\n", subkey.PublicKey.KeyId)
					}
				}
				keystr = ""
			} else {
				fmt.Printf("  %d) Invalid key: %g\n", i, err)
			}
		}
		if len(keyring) > 0 {
			err = nil
		}
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

func getFiles(walkdir, suffix string) []string {
	ret := []string{}
	err := filepath.Walk(walkdir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, suffix) {
				ret = append(ret, path)
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return ret
}
