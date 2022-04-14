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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	humanize "github.com/dustin/go-humanize"
)

var version = "test"
var debug, testOnly *bool
var attempts, shuffleAfter *int
var returnInt int
var getDisk, getMirror, getFails, getRecover int
var startTime = time.Now()
var uniqueCount, totalCount int
var totalBytes uint64
var duplicates *string
var after, before *time.Time
var logFile *io.Writer

var useList MirrorList // List of mirrors to use

type FileEntry struct {
	hash string
	size int
	path string
	dups []string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Yum Get RepoMD,  Version: %s\n\nUsage: %s [options...]\n\n", version, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Date formats supported: https://github.com/araddon/dateparse\n")
	}

	var mirrorList = flag.String("mirrors", "mirrorlist.txt", "Mirror / directory list of prefixes to use")
	var outputPath = flag.String("output", ".", "Path to put the repo files")
	var logFilePath = flag.String("log", "", "File in which to store a log of files downloaded")
	var threads = flag.Int("threads", 4, "Concurrent downloads")
	attempts = flag.Int("attempts", 40, "Attempts for each file")
	var connTimeout = flag.Duration("timeout", 10*time.Minute, "Max connection time, in case a mirror slows significantly")
	shuffleAfter = flag.Int("shuffle", 100, "Shuffle the mirror list every N downloads")
	var fileList = flag.String("list", "filelist.txt", "Filelist to be fetched (one per line with: HASH SIZE PATH)")
	debug = flag.Bool("debug", false, "Turn on debug comments")
	testOnly = flag.Bool("test", false, "Just validate downloaded files")
	duplicates = flag.String("dup", "symlink", "What to do with duplicates: omit, copy, symlink, hardlink")

	var afterStr = flag.String("after", "", "Select packages after specified date")
	var beforeStr = flag.String("before", "", "Select packages before specified date")

	flag.Parse()

	if *logFilePath != "" {
		logFile, err := os.Create(*logFilePath)
		if err != nil {
			log.Fatal("Error creating log file:", err)
		}
		defer logFile.Close()
	}

	if *afterStr != "" {
		t, err := dateparse.ParseLocal(*afterStr)
		after = &t
		if err != nil {
			log.Fatal("Error parsing after date", err)
		}
	}
	if *beforeStr != "" {
		t, err := dateparse.ParseLocal(*beforeStr)
		before = &t
		if err != nil {
			log.Fatal("Error parsing before date", err)
		}
	}

	switch *duplicates {
	case "omit", "copy", "symlink", "hardlink":
	default:
		fmt.Println("Invalid option for duplicates")
		return
	}

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
					total := getDisk + getMirror + getRecover
					percent := ""
					if uniqueCount > 0 {
						percent = fmt.Sprintf("%4.2f%%", 100.0*float32(total)/float32(uniqueCount))
					}
					fmt.Println("Stat:  OnDisk:", getDisk, "Downloaded:", getMirror, "Fails:",
						getFails, "Recovered:", getRecover, "Progress:", total, "/", uniqueCount, percent)
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
						Timeout: *connTimeout,
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
	jobs := make(chan *FileEntry, *threads)
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

	// Read line one by one and add them to the FileEntry map
	scanner := bufio.NewScanner(file)

	fileEntries := []*FileEntry{}
	fileMap := make(map[string]*FileEntry)
	for scanner.Scan() {
		line := scanner.Text()
		// Split line by line of filelist
		if strings.HasPrefix(line, "{") {
			parts := strings.SplitN(line, " ", 3)

			// If it is in an invalid format, skip
			if len(parts) < 3 {
				fmt.Println("Invalid format for input file list")
				continue
			}

			entry := FileEntry{}
			entry.size, err = strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("invalid file size", parts[1], "for", parts[2])
				continue
			}
			entry.hash = parts[0]
			entry.path = strings.TrimSpace(parts[2])
			mapstr := parts[1] + " " + parts[0]
			if e, ok := fileMap[mapstr]; ok {
				if e.path == entry.path {
					// They are the same path, ignore this line entirely
					continue
				}
				// We have been asked to not ignore duplicates
				if *duplicates != "omit" {
					if (isFile(e.path) && !isFile(entry.path)) ||
						(isSymlink(e.path) && !isSymlink(entry.path)) {
						// Flip the entries if one exists and the other doesn't, or
						// Keep to checking the linked-to-file instead of the link
						e.path, entry.path = entry.path, e.path
					}
					e.dups = append(e.dups, entry.path)
					if *duplicates == "copy" {
						totalBytes += uint64(entry.size)
					}
				}
			} else {
				// New entry, not already seen
				uniqueCount++
				totalBytes += uint64(entry.size)
				fileMap[mapstr] = &entry
				fileEntries = append(fileEntries, &entry)
			}
			totalCount++
		}
	}
	fmt.Println("# Files list stats -- Total:", totalCount, "Unique:", uniqueCount, "Size:", humanize.Bytes(totalBytes))
	for i, entry := range fileEntries {
		if i%(*shuffleAfter) == 0 {
			// Sort the mirror list by weight, latency + failures + random
			Shuffle()

			if *debug {
				fmt.Println("Mirror list:")
				useList.Print()
			}
		}
		jobs <- entry
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	close(jobs)
	for w := 1; w <= *threads; w++ {
		<-closure // ensure all the threads are closed
	}

	// Now that we have all the needed downloads, make all the links / copies
	for _, j := range fileEntries {
		for _, ln := range j.dups {
			from := path.Join(*outputPath, j.path)
			to := path.Join(*outputPath, ln)
			to_dir := filepath.Dir(to)
			rel_from, err := filepath.Rel(to_dir, from)
			//fmt.Println("from:", from, "to:", to, "rel:", rel_from)
			if err != nil {
				if *debug {
					fmt.Println("relative path", err)
				}
				continue
			}
			err = ensureDir(to_dir)
			if err != nil {
				if *debug {
					fmt.Println("ensure path", to_dir, "err:", err)
				}
				continue
			}
			switch *duplicates {
			case "copy":
				os.Remove(to)
				_, err = copyFile(from, to)
				if err != nil {
					if *debug {
						fmt.Println("err copying file", err)
					}
					continue
				}

			case "hardlink":
				os.Remove(to)
				err = os.Link(from, to)
				if err != nil {
					if *debug {
						fmt.Println("err making link", err)
					}
					continue
				}

			case "symlink":
				//if isSymlink(to) {
				os.Remove(to)
				//}
				err = os.Symlink(rel_from, to)
				if err != nil {
					if *debug {
						//fmt.Println("ln -s", rel_from, to)
						fmt.Println("err making link", err)
					}
					continue
				}
			}
		}
	}

	if !*testOnly {
		useList.Print()
		total := getDisk + getMirror + getRecover
		fmt.Println("Stat:  OnDisk:", getDisk, "Downloaded:", getMirror, "Fails:",
			getFails, "Recovered:", getRecover, "Progress:", total, "/", uniqueCount)
		if returnInt == 0 {
			fmt.Println("Successfully downloaded into", *outputPath)
		}
	} else {
		if returnInt == 0 {
			fmt.Println("# Successfully checked all files", *outputPath)
		}
	}
	os.Exit(returnInt)
}

var client = http.Client{
	Timeout: 5 * time.Second,
}

var worker_status []string

// This function is called when the downloads are threaded.  It's intended
// purpose is to loop over the FileEntry jobs which are sent over a channel to
// the threads whenever one opens up to grab the next job.  Once the jobs
// channel closes, the closure channel is used to make sure the threads are
// fully completed before continuing in the main function.
func worker(thread int, jobs <-chan *FileEntry, outputPath string, closure chan<- int) {
	for j := range jobs {
		//fmt.Printf("j = %+v\n", j)
		output := path.Join(outputPath, j.path)
		if *testOnly {
			err := handleFile(nil, j.hash, j.size, "", output)
			//fmt.Println("test", output)
			if err != nil {
				fmt.Printf("%s %d %s\n", j.hash, j.size, j.path)
				returnInt = 1
			}
			break
		}

		skip := []int{}
		isFail := false
		var url string
		var m *Mirror
		for len(skip) < *attempts && len(skip) < len(useList) {
			for m = PopWithout(skip); m == nil; m = PopWithout(skip) {
				if *debug {
					fmt.Println("  Waiting for a mirror to become available")
				}
				time.Sleep(3 * time.Second)
			}
			url = m.URL + "/" + strings.TrimPrefix(j.path, "/")
			worker_status[thread] = fmt.Sprintf("%d ~~ %s", m.ID, output)
			if *debug {
				fmt.Println("Downloading file", url, "to", output)
			}

			err := handleFile(m, j.hash, j.size, url, output)
			worker_status[thread] = "waiting"

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
			fmt.Println("  Exhausted the retries", j.path, skip)
			returnInt = 1
		}
	}
	closure <- 1
}

// The bulk of the work is done in this function, from testing the file on disk
// to see if it is valid, to downloading the file from a given mirror.
func handleFile(m *Mirror, hash string, size int, url, output string) error {
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
			return fmt.Errorf("File is directory")
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
				if *debug {
					fmt.Println("hash check failed", output, "hash:", hash, "err:", err)
				}
			}
		}
		if *testOnly {
			return fmt.Errorf("Invalid file")
		}
		if url == "" {
			return fmt.Errorf("url empty")
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

	fileTime, fileTimeErr := http.ParseTime(resp.Header.Get("Last-Modified"))
	if after != nil && fileTime.Before(*after) {
		return nil
	}
	if before != nil && fileTime.After(*before) {
		return nil
	}

	//file, err := os.Create(output)
	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
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
	file.Close()

	if fileSize != size {
		os.Remove(output)
		return fmt.Errorf("Size mismatch, %d != %s", fileSize, size)
	}

	// Check the hash and return any errors
	err = checkHash(hash, hashInterface)
	if err == nil {
		getMirror++
		if logFile != nil {
			fmt.Fprintln(*logFile, output)
		}
		success = true

		if fileTimeErr == nil {
			os.Chtimes(output, fileTime, fileTime)
		}
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
		return NewWithExpectedHash(strings.TrimPrefix(hash, "{alpine}"))
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
		if strings.TrimPrefix(hash, "{alpine}") == fmt.Sprintf("%s", h.Sum(nil)) {
			return nil
		}
	}
	return fmt.Errorf("Hash check failed")
}
