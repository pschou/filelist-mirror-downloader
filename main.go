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
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

var version = "test"
var debug, testOnly *bool
var attempts, shuffleAfter *int
var returnInt int
var getDisk, getMirror, getFails, getRecover int
var startTime = time.Now()

var useList MirrorList // List of mirrors to use

// HelloGet is an HTTP Cloud Function.
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Yum Get RepoMD,  Version: %s\n\nUsage: %s [options...]\n\n", version, os.Args[0])
		flag.PrintDefaults()
	}

	var mirrorList = flag.String("mirrors", "mirrorlist.txt", "Mirror / directory list of prefixes to use")
	var outputPath = flag.String("output", ".", "Path to put the repo files")
	var threads = flag.Int("threads", 4, "Concurrent downloads")
	attempts = flag.Int("attempts", 40, "Attempts for each file")
	var connTimeout = flag.Int("timeout", 300, "Max connection time, in case a mirror slows significantly")
	shuffleAfter = flag.Int("shuffle", 100, "Shuffle the mirror list ever N downloads")
	var fileList = flag.String("list", "filelist.txt", "Filelist to be fetched (one per line with: HASH SIZE PATH)")
	debug = flag.Bool("debug", false, "Turn on debug comments")
	testOnly = flag.Bool("test", false, "Just validate downloaded files")

	flag.Parse()

	if !*testOnly {
		fmt.Println("Press  [m] mirror list  [s] stats  [w] worker")
		go func() {
			// disable input buffering
			exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
			// do not display entered characters on the screen
			//exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
			var b []byte = make([]byte, 1)
			for {
				os.Stdin.Read(b)
				switch b[0] {
				case 'm':
					useList.Print()
				case 's':
					fmt.Println("Counts,  Disk:", getDisk, "Mirror:", getMirror, "Fails:", getFails, "Recovered:", getRecover)
				case 'w':
					if len(worker_status) > 0 {
						for i := 0; i < *threads; i++ {
							fmt.Printf(" %d) %s\n", i+1, worker_status[i])
						}
					}
				}
			}
		}()

		// Create the directory if needed
		err := ensureDir(*outputPath)
		if err != nil {
			log.Fatal(err)
		}

		mirrors := readMirrors(*mirrorList)
		var tmp string

		fmt.Println("Loaded", len(mirrors), "testing latencies and connectivity...")

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
						Timeout: time.Duration(*connTimeout) * time.Second,
					}).Dial,
					TLSHandshakeTimeout: 30 * time.Second,
				}

				useList = append(useList, Mirror{
					ID:      i + 1,
					URL:     m,
					Latency: delta,
					Client: http.Client{
						Timeout:   15 * time.Second,
						Transport: netTransport,
					},
				})
			}
		}
		fmt.Println("Downloading file list using", len(useList), "mirrors...")

		if len(useList) == 0 {
			log.Fatal("No mirrors found")
		}
	}

	// Setup a worker group to do work!
	jobs := make(chan string, *threads)
	closure := make(chan int, *threads)
	worker_status = make([]string, *threads)
	for w := 0; w < *threads; w++ {
		go worker(w, jobs, *outputPath, closure)
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
				Shuffle()

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

	if !*testOnly {
		useList.Print()
		fmt.Println("Counts,  Disk:", getDisk, "Mirror:", getMirror, "Fails:", getFails, "Recovered:", getRecover)
		if returnInt == 0 {
			fmt.Println("Successfully downloaded into", *outputPath)
		}
	}
	os.Exit(returnInt)
}

var client = http.Client{
	Timeout: 5 * time.Second,
}

var worker_status []string

func worker(thread int, jobs <-chan string, outputPath string, closure chan<- int) {
	for j := range jobs {
		parts := strings.SplitN(j, " ", 3)
		if len(parts) < 3 {
			fmt.Println("Invalid format for input file list")
		}
		size, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("invalid file size", parts[1], "for", parts[2])
		}
		output := path.Join(outputPath, parts[2])
		skip := []int{}
		isFail := false
		var url string
		var m *Mirror
		for len(skip) < *attempts {
			if !*testOnly {
				m = PopWithout(skip)
				if m == nil {
					//fmt.Println("  Exhausted the mirror list", parts[2], skip)
					returnInt = 1
					//useList.Print()
					break
				}
				url = m.URL + "/" + strings.TrimPrefix(parts[2], "/")
				worker_status[thread] = fmt.Sprintf("%d ~~ %s", m.ID, output)
				if *debug {
					fmt.Println("Downloading file", url, "to", output)
				}
			}

			err := downloadFile(m, parts[0], size, url, output)
			worker_status[thread] = "waiting"

			if *testOnly {
				if err != nil {
					fmt.Println(output)
					returnInt = 1
				}
				break
			}

			ClearUse(m.ID)
			if err != nil {
				skip = append(skip, m.ID)
				//fmt.Println("  Error:", err, "on", url, "using mirror id", m.ID)
				//fmt.Printf("  skip list: %+v\n", skip)
				//useList.Print()
				m.Failures++
				if !isFail {
					getFails++
				}
				isFail = true
				Shuffle()
			} else {
				if isFail {
					getRecover++
				}
				if *debug {
					fmt.Println("  Success on", url)
				}
				break
			}
		}
		if len(skip) == *attempts {
			fmt.Println("  Exhausted the retries", parts[2], skip)
			returnInt = 1
		}
	}
	closure <- 1
}

func downloadFile(m *Mirror, hash string, size int, url, output string) error {
	success := false
	defer func() {
		if !success {
			os.Remove(output)
		}
	}()

	dir, _ := path.Split(output)

	// Create the directory if needed
	err := ensureDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Check if file exists
	fileStat, err := os.Stat(output)

	// If we are in test mode, return if a file is missing
	if *testOnly && err != nil {
		return err
	}

	if err == nil {
		if fileStat.IsDir() {
			fmt.Println("  File is directory", output)
			return err
		}
		if int(fileStat.Size()) != size {
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
			hashInterface := getHash(hash)
			io.Copy(hashInterface, file)
			file.Close()

			// Check the hash and return any errors
			err = checkHash(hash, hashInterface)
			if err == nil {
				if *debug {
					fmt.Println("  Skipping, found valid file:", output)
				}
				getDisk++
				success = true
				return nil
			} else {
				fmt.Println("hash check failed", output, "hash:", hash, "err:", err)
			}
		}
		if url == "" {
			return err
		}
		os.Remove(output)
	}

	hashInterface := getHash(hash)
	if hashInterface == nil {
		return fmt.Errorf("Unknown hash type for: %s", hash)
	}

	//resp, err := m.Client.Get(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error in HTTP making new request", err)
		return err
	}

	ctx, cancel := context.WithTimeout(req.Context(), 60*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	start := time.Now()
	resp, err := m.Client.Do(req)

	if err != nil {
		log.Println("Error in HTTP get request", err)
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(output)
	if err != nil {
		return err
	}

	buf := make([]byte, 10000)
	var readBytes, fileSize int
	var readErr error

	var zeroCounter int

	// Do the download!
	for readErr != io.EOF {
		// read from webserver
		readBytes, readErr = resp.Body.Read(buf)
		m.Bytes += readBytes
		tick := time.Now()
		m.Time += tick.Sub(start)
		start = tick
		if fileSize+readBytes > size {
			readBytes = size - fileSize
			readErr = io.EOF
		}
		fileSize += readBytes

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

		_, err = file.Write(buf[:readBytes])
		if err != nil {
			return err
		}

		_, err = hashInterface.Write(buf[:readBytes])
		if err != nil {
			return err
		}
	}

	if fileSize != size {
		os.Remove(output)
		return fmt.Errorf("Size mismatch, %d != %s", fileSize, size)
	}

	// Check the hash and return any errors
	err = checkHash(hash, hashInterface)
	if err == nil {
		getMirror++
		success = true
	}
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
		//fmt.Println("fn hash check", strings.TrimPrefix(hash, "{alpine}"), fmt.Sprintf("%s", h.Sum(nil)))
		if strings.EqualFold(strings.TrimPrefix(hash, "{alpine}"), fmt.Sprintf("%s", h.Sum(nil))) {
			return nil
		}
	}
	return fmt.Errorf("Hash check failed")
}
