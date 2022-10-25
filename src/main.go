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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/araddon/dateparse"
	humanize "github.com/dustin/go-humanize"
	tease "github.com/pschou/go-tease"
)

var version = "test"
var debug *bool

//var debug, testOnly *bool
var attempts, shuffleAfter *int
var returnInt int
var getDisk, getDownloaded, getUncomp, getFails, getRecover, getSkip int
var startTime = time.Now()
var uniqueCount, totalCount int
var totalBytes uint64
var duplicates *string
var after, before *time.Time
var logFile *os.File
var outputPath, keyringFile *string

var useListMutex sync.Mutex
var useList MirrorList // List of mirrors to use

type FileEntry struct {
	hash                string
	size                int
	path                string
	ext                 string
	compressedVersion   []*FileEntry
	uncompressedVersion []*FileEntry
	modified            time.Time
	dups                []string
	attempted           bool
	success             bool
	skip                []int
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Yum Get RepoMD,  Version: %s\n\nUsage: %s [options...]\n\n", version, os.Args[0])
		flag.PrintDefaults()
		//fmt.Fprintf(os.Stderr, "Date formats supported: https://github.com/araddon/dateparse\n")
	}

	var background = flag.Bool("background", false, "Ignore all keyboard inputs, background mode")
	var mirrorList = flag.String("mirrors", "mirrorlist.txt", "Mirror / directory list of prefixes to use")
	outputPath = flag.String("output", ".", "Path to put the repo files")
	var logFilePath = flag.String("log", "", "File in which to store a log of files downloaded\n"+
		"Line prefixes (OnDisk, OnDiskSkip, Skipped, Fetched, Uncompressed, Failed), indicate action taken.\n"+
		"Skip means that a file falls outside the required date bounds")
	var threads = flag.Int("threads", 4, "Concurrent downloads")
	attempts = flag.Int("attempts", 40, "Attempts for each file")
	var connTimeout = flag.Duration("timeout", 10*time.Minute, "Max connection time per file, in case a mirror slows significantly\n"+
		"If one is downloading large ISO files, a longer time may be needed.")
	shuffleAfter = flag.Int("shuffle", 100, "Shuffle the mirror list every N downloads")
	var fileList = flag.String("list", "filelist.txt", "Filelist to be fetched (one per line with: HASH SIZE PATH)")
	debug = flag.Bool("debug", false, "Turn on debug comments")
	//testOnly = flag.Bool("test", false, "Just validate downloaded files")
	duplicates = flag.String("dup", "symlink", "What to do with duplicates: omit, copy, symlink, hardlink")
	keyringFile = flag.String("keyring", "", "Use keyring for verifying signed package files (example: keyring.gpg or keys/ directory)")

	var afterStr = flag.String("after", "", "Select packages after specified date\n"+
		"Date formats supported: https://github.com/araddon/dateparse")
	var beforeStr = flag.String("before", "", "Select packages before specified date\n"+
		"Date formats supported: https://github.com/araddon/dateparse")

	flag.Parse()

	if *keyringFile != "" {
		fmt.Println("Loading keys from", *keyringFile)
		if _, ok := isDirectory(*keyringFile); ok {
			//keyring = openpgp.EntityList{}
			for _, file := range getFiles(*keyringFile, []string{".pub", ".gpg", ".pgp"}) {
				fmt.Println("loading key", file)
				gpgFile, err := os.ReadFile(file)
				if err != nil {
					log.Fatal("Error reading keyring file", err)
				}
				p, r, err := loadKeys(gpgFile)
				if err != nil {
					log.Fatal("Error loading keyring file", err)
				}
				//fmt.Println("  found", len(p)+len(r), "keys")
				//if pgpkey, ok := fileKeys.(openpgp.EntityList); ok {
				pgpKeys = append(pgpKeys, p...)
				rsaKeys = append(rsaKeys, r...)
				//}
			}
		} else {
			gpgFile, err := os.ReadFile(*keyringFile)
			if err != nil {
				log.Fatal("Error reading keyring file", err)
			}
			p, r, err := loadKeys(gpgFile)
			if err != nil {
				log.Fatal("Error loading keyring file", err)
			}
			//if pgpkey, ok := fileKEys.(openpgp.EntityList); ok {
			pgpKeys = append(pgpKeys, p...)
			rsaKeys = append(rsaKeys, r...)
			//}
		}
		fmt.Println("Key count:", len(pgpKeys), "pgp keys", len(rsaKeys), "rsa keys")
	}

	if *logFilePath != "" {
		var err error
		logFile, err = os.Create(*logFilePath)
		if err != nil {
			log.Fatal("Error creating log file:", err)
		}
		defer logFile.Sync()
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

	if !*background {
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
					total := getDisk + getDownloaded + getUncomp
					percent := ""
					if uniqueCount > 0 {
						percent = fmt.Sprintf("%4.2f%%", 100.0*float32(total)/float32(uniqueCount))
					}
					fmt.Println("Stat:  OnDisk:", getDisk, "Downloaded:", getDownloaded, "Uncompressed:", getUncomp, "Fails:",
						getFails, "Skipped:", getSkip, "Recovered:", getRecover, "Progress:", total, "/", uniqueCount, percent)
				case 'w':
					if len(worker_status) > 0 {
						for i := 0; i < *threads; i++ {
							fmt.Printf(" %d) %s\n", i+1, worker_status[i])
						}
					}
				}
			}
		}()
	}

	// Create the directory if needed
	err := ensureDir(*outputPath)
	if err != nil {
		log.Fatal(err)
	}

	mirrors := readMirrors(*mirrorList)
	var tmp string

	fmt.Println("Loaded", len(mirrors), "testing latencies and connectivity...")

	var wg sync.WaitGroup

	// Test speeds
	id := 0
	for _, mm := range mirrors {
		mirror_url, err := url.Parse(mm)
		if err != nil {
			continue
		}
		if *debug {
			fmt.Println("looking up ip for:", mirror_url.Host)
		}
		host, _, err := net.SplitHostPort(mirror_url.Host)
		if err != nil {
			host = mirror_url.Host
		}

		ips := getIPs(host)
		for _, ip := range ips {
			if *debug {
				fmt.Println("  ip:", ip)
			}

			wg.Add(1)

			go func(m string, ip net.IP) {
				defer wg.Done()
				if *debug {
					fmt.Println("Starting test on", m)
				}
				//repoPath := m + "/" + repoPath + "/"
				//repomdPath := repoPath + "repodata/repomd.xml"

				dial := func(network, address string) (net.Conn, error) {
					_, port, err := net.SplitHostPort(address)
					if err != nil {
						return nil, err
					}
					return (&net.Dialer{
						Timeout: *connTimeout,
					}).Dial(network, net.JoinHostPort(ip.String(), port))
				}

				var netTransport = &http.Transport{
					Dial:                dial,
					TLSHandshakeTimeout: 30 * time.Second,
				}
				client := http.Client{
					Timeout:   5 * time.Second,
					Transport: netTransport,
				}

				start := time.Now()
				tmp = readFile(m, client)
				if len(tmp) < 100 {
					return
				}
				delta := time.Now().Sub(start).Seconds() * 1000
				if *debug {
					fmt.Printf("  %.02f ms for %d bytes - %s\n", delta, len(tmp), m)
				}
				if delta < 2000 {

					title := m
					if len(ips) > 1 {
						title += " (" + ip.String() + ")"
					}
					client.Timeout = *connTimeout

					useListMutex.Lock()
					id++
					useList = append(useList, Mirror{
						ID:      id,
						title:   title,
						URL:     m,
						IP:      ip,
						Latency: delta,
						Client:  client,
					})
					useListMutex.Unlock()
				}
			}(mm, ip)
			time.Sleep(70 * time.Millisecond)
		}
	}
	wg.Wait()

	// Expand the mirror list to allow the threading to work
	for len(useList) < *threads {
		count := len(useList)
		for _, ul := range useList {
			ul.ID += count
			useList = append(useList, ul)
		}
	}

	Shuffle()
	fmt.Println("Downloading file list using", len(useList), "mirrors...")

	if len(useList) == 0 {
		log.Fatal("No mirrors found")
	}

	// Setup a worker group to do work!
	jobs := make(chan *FileEntry, *threads)
	worker_status = make([]string, *threads)

	for w := 0; w < *threads; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
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
	compressed := []*FileEntry{}
	uncompressed := make(map[string]*FileEntry)
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

			entry := FileEntry{
				hash: parts[0],
				path: strings.TrimSpace(parts[2]),
				ext:  filepath.Ext(strings.TrimSpace(parts[2])),
			}
			entry.size, err = strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("invalid file size", parts[1], "for", parts[2])
				continue
			}
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
				// New entry, not already matched

				switch entry.ext {
				case ".gz", ".xz":
					// add to compressed file list if compressed file extension
					compressed = append(compressed, &entry)
				default:
					uncompressed[entry.path] = &entry
				}

				uniqueCount++
				totalBytes += uint64(entry.size)
				fileMap[mapstr] = &entry
				fileEntries = append(fileEntries, &entry)
			}
			totalCount++
		}
	}

	// Loop over compressed files and add compressed to uncompressed
	for _, entry := range compressed {
		for _, entry_path := range append(entry.dups, entry.path) {
			if len(entry_path) > len(entry.ext) {
				uc_name := entry_path[:len(entry_path)-len(entry.ext)]
				if uc_entry, ok := uncompressed[uc_name]; ok {
					// Found compressed version of file
					uc_entry.compressedVersion = append(uc_entry.compressedVersion, entry)
					entry.uncompressedVersion = append(entry.uncompressedVersion, uc_entry)
				}
			}
		}
	}

	fmt.Println("# Files list stats -- Total:", totalCount, "Unique:", uniqueCount, "Size:", humanize.Bytes(totalBytes))
	for i, entry := range fileEntries {
		if len(entry.compressedVersion) > 0 {
			// Skip downloading entries with available compressed version of file to reduce network transfers, the compressed version
			// should be tried first.
			continue
		}
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
	if *debug {
		fmt.Print("Waiting for threads to close")
	}
	wg.Wait()

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

	for _, j := range fileEntries {
		if !j.success {
			if *debug {
				fmt.Println("Failed:", j.path)
			}
			if logFile != nil {
				output := path.Join(*outputPath, j.path)
				fmt.Fprintln(logFile, "Failed:", output)
			}
			returnInt = 1
		}
	}

	useList.Print()
	total := getDisk + getDownloaded + getUncomp
	fmt.Println("Stat:  OnDisk:", getDisk, "Downloaded:", getDownloaded, "Uncompressed:", getUncomp, "Fails:",
		getFails, "Skipped:", getSkip, "Recovered:", getRecover, "Progress:", total, "/", uniqueCount)
	if returnInt == 0 {
		fmt.Println("Successfully downloaded into", *outputPath)
	}
	os.Exit(returnInt)
}

var worker_status []string

// This function is called when the downloads are threaded.  It's intended
// purpose is to loop over the FileEntry jobs which are sent over a channel to
// the threads whenever one opens up to grab the next job.  Once the jobs
// channel closes, the closure channel is used to make sure the threads are
// fully completed before continuing in the main function.
func worker(thread int, jobs <-chan *FileEntry, wg *sync.WaitGroup) {
	defer func() {
		if *debug {
			fmt.Println("Closing thread", thread)
		}
		worker_status[thread] = "closed"
		wg.Done()
	}()
	worker_status[thread] = "init"
	for next_job := range jobs {
		worker_status[thread] = fmt.Sprintf("finding mirror for %s", next_job.path)

		func() {
			m, any_left := GetMirrorOrQueue(next_job)
			if m != nil {
				defer ClearUse(m.ID)
				for _, j := range append(GetQueue(m.ID), next_job) {
					process(thread, m, j)
				}
			}
			if !any_left {
				fmt.Println("Failed: ", next_job.path)
				getFails++
			}
		}()

		process_straggler := func() bool {
			jobs, m := FindStragglers()
			if m != nil {
				defer ClearUse(m.ID)
			} else {
				return false
			}
			for _, j := range jobs {
				process(thread, m, j)
			}
			return true
		}

		for process_straggler() {
		}

		worker_status[thread] = "waiting"
	}
}

func process(thread int, m *Mirror, j *FileEntry) {
	attempted := j.attempted
	j.attempted = true
	output := path.Join(*outputPath, j.path)
	url := m.URL + "/" + strings.TrimPrefix(j.path, "/")
	worker_status[thread] = fmt.Sprintf("%d ~~ %s", m.ID, output)
	if *debug {
		fmt.Println("Downloading file", url, "to", output)
	}

	err := handleFile(m, j)
	if err != nil {
		m.Failures++
		if *debug {
			fmt.Printf("Failed file: %s %d %s %s\n", j.hash, j.size, j.path, err)
		}
	}
	if err == nil {
		j.success = true
		if attempted {
			getRecover++
		}

		for _, uj := range j.uncompressedVersion {
			switch j.ext {
			case ".gz", ".xz":
				hashInterface := getHash(uj.hash)
				if hashInterface == nil {
					log.Println("Unknown hash interface:", uj.hash)
					continue
				}
				uj_output := path.Join(*outputPath, uj.path)
				_, err = copyUncompFile(output, uj_output, j.ext, hashInterface)
				if err != nil {
					log.Println("Error in uncompress:", err)
					os.Remove(uj_output)
					Queue(uj)
					continue
				}

				// Check the hash and return any errors
				err = checkHash(uj.hash, hashInterface)
				if err == nil {
					getUncomp++
					fmt.Fprintln(logFile, "Uncompressed:", uj_output)
					uj.success = true
					if &j.modified != nil {
						os.Chtimes(uj_output, j.modified, j.modified)
					}
				} else {
					log.Println("Error in uncompress hash check:", err)
					os.Remove(uj_output)
					Queue(uj)
				}
			default:
				log.Println("Unsupported compression type", j.ext)
			}
		}
	} else {
		j.skip = append(j.skip, m.ID)
		hasQueued := Queue(j)
		if !hasQueued {
			// Could not re-queue, count as a fail
			getFails++
		} else if *debug {
			fmt.Printf("requeued %s\n", j.path)
		}
	}
}

// The bulk of the work is done in this function, from testing the file on disk
// to see if it is valid, to downloading the file from a given mirror.
func handleFile(m *Mirror, j *FileEntry) error {
	output := path.Join(*outputPath, j.path)
	url := m.URL + "/" + strings.TrimPrefix(j.path, "/")
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
	/*if *testOnly && err != nil {
		return err
	}*/

	if err == nil {
		if fileStat.IsDir() {
			fmt.Println("  File is directory", output)
			return fmt.Errorf("File is directory")
		}
		if int(fileStat.Size()) != j.size {
			if *debug {
				fmt.Println("  Mismatched size of file", output, int(fileStat.Size()), "!=", j.size)
			}
		} else {
			file, err := os.Open(output)
			if err != nil {
				return err
			}
			if *debug {
				fmt.Println("  Found file:", output)
			}
			hashInterface := getHash(j.hash)
			if hashInterface == nil {
				return fmt.Errorf("Unknown hash type: %q", j.hash)
			}
			io.Copy(hashInterface, file)
			file.Close()

			// Check the hash and return any errors
			err = checkHash(j.hash, hashInterface)
			if err == nil {
				if *debug {
					fmt.Println("  Skipping, found valid file:", output)
				}
				getDisk++

				if logFile != nil {
					stat, err := os.Stat(output)
					if err == nil && ((after != nil && stat.ModTime().Before(*after)) || (before != nil && stat.ModTime().After(*before))) {
						fmt.Fprintln(logFile, "OnDiskSkip:", output)
					} else {
						fmt.Fprintln(logFile, "OnDisk:", output)
					}
				}

				success = true
				return nil
			} else {
				if *debug {
					fmt.Println("hash check failed", output, "hash:", j.hash, "err:", err)
				}
			}
		}
		/*if *testOnly {
			return fmt.Errorf("Invalid file")
		}*/
		if url == "" {
			return fmt.Errorf("url empty")
		}
		os.Remove(output)
	}

	// Build a hash interface for verification of downloads, reading bytes 0-N
	hashInterface := getHash(j.hash)
	if hashInterface == nil {
		return fmt.Errorf("Unknown hash type for: %s", j.hash)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error in HTTP making new request", err)
		return err
	}

	req.Header.Set("User-Agent", "curl/7.29.0")
	//ctx, cancel := context.WithTimeout(req.Context(), 120*time.Second)
	//defer cancel()

	//req = req.WithContext(ctx)

	start := time.Now()
	resp, err := m.Client.Do(req)

	if err != nil {
		log.Println("Error in HTTP get request", err)
		return err
	}
	defer resp.Body.Close()

	fileTime, fileTimeErr := http.ParseTime(resp.Header.Get("Last-Modified"))
	if fileTimeErr == nil {
		j.modified = fileTime
	}

	respBody := io.Reader(resp.Body)

	// Tease out reading the file to get the header for the date information
	if fp, ok := file_parser[filepath.Ext(output)]; ok {
		tr := tease.NewReader(resp.Body)
		fileDetail := fp(tr)
		tr.Seek(0, io.SeekStart) // Return to the beginning of the tease reader
		tr.Pipe()                // Flatten the tease reader to an io.Reader
		respBody = tr            // Set the respBody to the tease reader
		if fileDetail != nil {
			if *debug {
				fmt.Println("Parsing time for", output, "-- server:", fileTime.UTC(), "file:", fileDetail.time.UTC())
			}
			fileTime = fileDetail.time
		}
	}

	// If the after time is set and the time is before the after time, skip the file
	if after != nil && fileTime.Before(*after) {
		if logFile != nil {
			fmt.Fprintln(logFile, "Skipped:", output)
		}
		getSkip++
		return nil
	}

	// If the before time is set and the time is after the before time, skip the file
	if before != nil && fileTime.After(*before) {
		if logFile != nil {
			fmt.Fprintln(logFile, "Skipped:", output)
		}
		getSkip++
		return nil
	}

	// Open the file for writing
	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}

	targetWriters := []io.Writer{file, hashInterface}

	var checksumResult = make(chan bool)
	var pipeWriter *io.PipeWriter

	// If a keyring has been specified, include a GPG check
	if *keyringFile != "" {
		// If we don't have a handler for a file type, we fail to letting it pass,
		// as this is a secondary check, the checksum is the primary check
		if fp, ok := file_sig_check[filepath.Ext(output)]; ok {
			var pipeReader *io.PipeReader
			pipeReader, pipeWriter = io.Pipe()
			targetWriters = append(targetWriters, pipeWriter)
			go fp(pipeReader, output, checksumResult)
		}
	}

	writer := io.MultiWriter(targetWriters...)

	buf := make([]byte, 10000)
	var readBytes, fileSize int
	var readErr error

	var zeroCounter int

	// Do the download!
	for readErr != io.EOF {
		// read from webserver
		readBytes, readErr = respBody.Read(buf)

		{ // Track download speed
			m.Bytes += readBytes
			tick := time.Now()
			m.Time += tick.Sub(start)
			start = tick
		}

		if fileSize+readBytes > j.size {
			readBytes = j.size - fileSize
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
		if err != nil {
			return err
		}

		_, err = writer.Write(buf[:readBytes])
		if err != nil {
			return err
		}
	}
	file.Close()

	if fileSize != j.size {
		os.Remove(output)
		return fmt.Errorf("Size mismatch, %d != %d", fileSize, j.size)
	}

	// Check the hash and return any errors
	err = checkHash(j.hash, hashInterface)
	if err == nil {
		getDownloaded++

		// Now we can also verify the key signature if the check was started
		if pipeWriter != nil {
			pipeWriter.Close()
			if <-checksumResult {
				if *debug {
					fmt.Println("Signature check passed or skipped")
				}
			} else {
				if *debug {
					fmt.Println("Signature check failed")
				}
				if logFile != nil {
					fmt.Fprintln(logFile, "SignatureFail:", output)
				}
				// We return nil here as the file passed the checksum but failed
				// signature check, it's a out-of-trust signed file in a repo, and
				// like a bad apple in the bushel, should not be trusted.
				return nil
			}
		}

		if logFile != nil {
			fmt.Fprintln(logFile, "Fetched:", output)
		}
		success = true

		if fileTimeErr == nil {
			os.Chtimes(output, fileTime, fileTime)
		}
	}
	return err
}

func getHash(hash string) hash.Hash {
	hash = strings.ToLower(hash)
	switch {
	case strings.HasPrefix(hash, "{none}"):
		return sha1.New()
	case strings.HasPrefix(hash, "{sha1}"):
		return sha1.New()
	case strings.HasPrefix(hash, "{sha}"):
		return sha1.New()
	case strings.HasPrefix(hash, "{md5}"):
		return md5.New()
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

func checkHash(orig_hash string, h hash.Hash) error {
	hash := strings.ToLower(orig_hash)
	switch {
	case strings.HasPrefix(hash, "{none}"):
		return nil
	case strings.HasPrefix(hash, "{sha1}"):
		if strings.EqualFold(strings.TrimPrefix(hash, "{sha1}"), fmt.Sprintf("%x", h.Sum(nil))) {
			return nil
		}
	case strings.HasPrefix(hash, "{sha}"):
		if strings.EqualFold(strings.TrimPrefix(hash, "{sha}"), fmt.Sprintf("%x", h.Sum(nil))) {
			return nil
		}
	case strings.HasPrefix(hash, "{md5}"):
		if strings.EqualFold(strings.TrimPrefix(hash, "{md5}"), fmt.Sprintf("%x", h.Sum(nil))) {
			return nil
		}
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
		if orig_hash[len("{alpine}"):] == fmt.Sprintf("%s", h.Sum(nil)) {
			return nil
		}
	}
	return fmt.Errorf("Hash check failed")
}
