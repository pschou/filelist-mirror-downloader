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
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ulikunitz/xz"
)

var ensureDirMap = make(map[string]int8)
var ensureDirSync sync.Mutex

func ensureDir(dirName string) error {
	ensureDirSync.Lock()
	defer ensureDirSync.Unlock()

	if _, ok := ensureDirMap[dirName]; ok {
		return nil
	}
	ensureDirMap[dirName] = 1

	err := os.MkdirAll(dirName, 0755)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

func isFile(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func isSymlink(filepath string) bool {
	fi, err := os.Lstat(filepath)
	if err != nil {
		return false
	}
	// ..check err...
	return fi.Mode()&os.ModeSymlink == os.ModeSymlink
	// This is a symlink
}

func readFile(filePath string, client http.Client) string {
	// Declare file handle for the reading
	var file io.Reader

	if _, err := os.Stat(filePath); err == nil {
		log.Println("Reading in file", filePath)

		// Open our xmlFile
		rawFile, err := os.Open(filePath)
		if err != nil {
			//log.Println("Error in HTTP get request", err)
			return ""
		}

		// Make sure the file is closed at the end of the function
		defer rawFile.Close()
		file = rawFile
	} else if strings.HasPrefix(filePath, "http") {

		req, err := http.NewRequest("GET", filePath, nil)
		if err != nil {
			log.Println("Error in building HTTP request", err)
			return ""
		}

		req.Header.Set("User-Agent", "curl/7.29.0")
		resp, err := client.Do(req)
		if err != nil {
			//log.Println("Error in HTTP get request", err)
			return ""
		}

		defer resp.Body.Close()
		file = resp.Body
	} else if _, err := os.Stat(filePath); err == nil {
		log.Println("Reading in file", filePath)

		// Open our xmlFile
		rawFile, err := os.Open(filePath)
		if err != nil {
			log.Println("Error opening file locally", err)
			return ""
		}

		// Make sure the file is closed at the end of the function
		defer rawFile.Close()
		file = rawFile
	} else {
		log.Println("Unable to open file", filePath)
		return ""
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	return buf.String()
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func copyUncompFile(src, dst, ext string, hash io.Writer) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	var uncomp_src io.Reader
	switch ext {
	case ".gz":
		gunzip, err := gzip.NewReader(source)
		if err != nil {
			return 0, err
		}
		defer gunzip.Close()
		uncomp_src = gunzip
	case ".xz":
		xz_unzip, err := xz.NewReader(source)
		if err != nil {
			return 0, err
		}
		uncomp_src = xz_unzip
	}

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(io.MultiWriter(destination, hash), uncomp_src)
	return nBytes, err
}
