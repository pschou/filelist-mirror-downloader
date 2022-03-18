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
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var version = "test"
var debug *bool
var attempts, shuffleAfter *int
var returnInt int

// HelloGet is an HTTP Cloud Function.
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Yum Get RepoMD,  Version: %s\n\nUsage: %s [options...]\n\n", version, os.Args[0])
		flag.PrintDefaults()
	}

	var mirrorList = flag.String("mirrors", "mirrorlist.txt", "Mirror / directory list of prefixes to use")
	var outputPath = flag.String("output", ".", "Path to put the repo files")
	var threads = flag.Int("threads", 4, "Concurrent downloads")
	attempts = flag.Int("attempts", 10, "Attempts for each file")
	shuffleAfter = flag.Int("shuffle", 10, "Shuffle the mirror list ever N downloads")
	var fileList = flag.String("list", "filelist.txt", "Filelist to be fetched (one per line with: HASH SIZE PATH)")
	debug = flag.Bool("debug", false, "Turn on debug comments")

	flag.Parse()

	// Create the directory if needed
	err := ensureDir(*outputPath)
	if err != nil {
		log.Fatal(err)
	}

	mirrors := readMirrors(*mirrorList)
	var tmp string
	var useList MirrorList // List of mirrors to use

	// Test speeds
	for i, m := range mirrors {
		//repoPath := m + "/" + repoPath + "/"
		//repomdPath := repoPath + "repodata/repomd.xml"
		start := time.Now()
		tmp = readFile(m)
		delta := time.Now().Sub(start).Seconds() * 1000
		if *debug {
			fmt.Printf("%d) %.02f ms for %d bytes - %s\n", i, delta, len(tmp), m)
		}
		if delta < 4000 && len(tmp) > 100 {

			var netTransport = &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 15 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 15 * time.Second,
			}

			useList = append(useList, Mirror{
				ID:      i + 1,
				URL:     m,
				Latency: delta,
				Client: http.Client{
					Timeout:   5 * time.Second,
					Transport: netTransport,
				},
			})
		}
	}

	if len(useList) == 0 {
		log.Fatal("No mirrors found")
	}

	// Setup a worker group to do work!
	jobs := make(chan string, *threads)
	closure := make(chan int, *threads)
	for w := 1; w <= *threads; w++ {
		go worker(useList, jobs, *outputPath, closure)
	}

	file, err := os.Open(*fileList)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read line one by one and send to threads
	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "{") {
			i++
			if i%(*shuffleAfter) == 0 {
				// Sort the mirror list by weight, latency + failures + random
				useList.Shuffle()

				if *debug {
					fmt.Println("Mirror list:")
					useList.Print()
				}
			}
			jobs <- scanner.Text()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	close(jobs)
	for w := 1; w <= *threads; w++ {
		<-closure // ensure all the threads are closed
	}

	os.Exit(returnInt)
}

var client = http.Client{
	Timeout: 5 * time.Second,
}

func worker(mirrors MirrorList, jobs <-chan string, outputPath string, closure chan<- int) {
	for j := range jobs {
		parts := strings.SplitN(j, " ", 3)
		if len(parts) < 3 {
			fmt.Println("Invalid format for input file list")
		}
		output := path.Join(outputPath, parts[2])
		skip := []int{}
		for len(skip) < *attempts {
			m := mirrors.PopWithout(skip)
			if m == nil {
				fmt.Println("  Exhausted the mirror list, no additional mirrors to try")
				break
			}
			url := m.URL + "/" + strings.TrimPrefix(parts[2], "/")
			if *debug {
				fmt.Println("Downloading file", url, "to", output)
			}
			err := downloadFile(m, parts[0], parts[1], url, output)
			if err != nil {
				skip = append(skip, m.ID)
				fmt.Println("  Error:", err)
				m.Failures++
				m.mirrors.Shuffle()
			} else {
				if *debug {
					fmt.Println("  Success on", url)
				}
				break
			}
		}
		if len(skip) == *attempts {
			returnInt = 1
		}
	}
	closure <- 1
}

func downloadFile(m *Mirror, hash, size, url, output string) error {
	success := false
	defer func() {
		m.InUse = false
		if !success {
			os.Remove(output)
		}
	}()
	//resp, err := m.Client.Get(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error in HTTP making new request", err)
		return err
	}

	ctx, cancel := context.WithTimeout(req.Context(), 15*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := m.Client.Do(req)

	if err != nil {
		log.Println("Error in HTTP get request", err)
		return err
	}
	defer resp.Body.Close()

	dir, _ := path.Split(output)
	hashInterface := getHash(hash)

	// Create the directory if needed
	err = ensureDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Check if file exists
	fileStat, err := os.Stat(output)
	if err == nil {
		if fileStat.IsDir() {
			fmt.Println("  File is directory", output)
			log.Fatal(err)
		}
		if fmt.Sprintf("%d", fileStat.Size()) != size {
			if *debug {
				fmt.Println("  Mismatched size of file", output)
			}
		} else {
			file, err := os.Open(output)
			if err != nil {
				return err
			}
			if *debug {
				fmt.Println("  Found file:", output)
			}
			io.Copy(hashInterface, file)
			if checkHash(hash, hashInterface) == nil {
				if *debug {
					fmt.Println("  Skipping, found valid file:", output)
				}
				success = true
				return nil
			}
			hashInterface = getHash(hash)
		}
		os.Remove(output)
	}

	file, err := os.Create(output)
	if err != nil {
		return err
	}

	if hashInterface == nil {
		return fmt.Errorf("Unknown hash type for: %s", hash)
	}
	buf := make([]byte, 10000)
	var readBytes, fileSize int
	var readErr error

	var zeroCounter int

	// Do the download!
	for readErr != io.EOF {
		// read from webserver
		readBytes, readErr = resp.Body.Read(buf)
		if readErr != nil && readErr != io.EOF && *debug {
			fmt.Println("  Error on reading:", readErr, "with bytes", readBytes)
		}

		if readBytes == 0 {
			if readErr != nil {
				return readErr
			}
			zeroCounter++
			if zeroCounter > 1000 {
				return fmt.Errorf("Server stopped talking: %s", url)
			}
		} else {
			zeroCounter = 0
		}
		//log.Println("  Read in", readBytes)
		if err != nil {
			return err
		}
		fileSize += readBytes

		_, err = file.Write(buf[:readBytes])
		if err != nil {
			return err
		}

		_, err = hashInterface.Write(buf[:readBytes])
		if err != nil {
			return err
		}
	}

	if fmt.Sprintf("%d", fileSize) != size {
		os.Remove(output)
		return fmt.Errorf("Size mismatch, %d != %s", fileSize, size)
	}

	// Check the hash and return any errors
	err = checkHash(hash, hashInterface)
	if err != nil {
		os.Remove(output)
	}
	success = true
	return nil
}

func getHash(hash string) hash.Hash {
	switch {
	case strings.HasPrefix(hash, "{sha256}"):
		return sha256.New()
	case strings.HasPrefix(hash, "{sha512}"):
		return sha512.New()
	case strings.HasPrefix(hash, "{alpine}"):
		return NewWithExpectedHash(hash[8:])
	default:
		if *debug {
			fmt.Println("Unknown hash type:", hash)
		}
	}
	return nil
}

func checkHash(hash string, h hash.Hash) error {
	switch {
	case strings.HasPrefix(hash, "{sha256}"):
		if strings.EqualFold(strings.TrimPrefix(hash, "{sha256}"), fmt.Sprintf("%x", h.Sum(nil))) {
			return nil
		}
	case strings.HasPrefix(hash, "{sha512}"):
		if strings.EqualFold(strings.TrimPrefix(hash, "{sha512}"), fmt.Sprintf("%x", h.Sum(nil))) {
			return nil
		}
	case strings.HasPrefix(hash, "{alpine}"):
		if strings.EqualFold(strings.TrimPrefix(hash, "{alpine}"), fmt.Sprintf("%s", h.Sum(nil))) {
			return nil
		}
	}
	return fmt.Errorf("Hash check failed")
}
