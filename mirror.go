package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Mirror struct {
	ID       int
	URL      string
	Latency  float64
	Random   float64
	Failures int
	Client   http.Client
	Bytes    int
	Time     time.Duration
	inUse    bool
	queue    []*FileEntry
	c        chan struct{}
	c_n      int
}

func readMirrors(mirrorFile string) []string {
	file, err := os.Open(mirrorFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var line string

	ret := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		line = strings.TrimSuffix(line, "/")
		ret = append(ret, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ret
}

type MirrorList []Mirror

var MirrorListSync sync.Mutex

func (m MirrorList) Len() int { return len(m) }
func (m MirrorList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m MirrorList) Less(i, j int) bool {
	return m[i].Latency+float64(m[i].Failures)*20+m[i].Random < m[j].Latency+float64(m[j].Failures)*20+m[j].Random
}
func (m MirrorList) Print() {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	fmt.Println(" #  Weight Latency Fails   Rand InUse MBps URL")
	for _, e := range m {
		fmt.Printf("%2d) %6.02f %6.02f %6d %6.02f %t %5.1f %s\n", e.ID, e.Latency+float64(e.Failures)*20+e.Random, e.Latency, e.Failures, e.Random, e.inUse, float32(e.Bytes)/float32(e.Time)*1e3, e.URL)
	}
}
func Shuffle() {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	for i := range useList {
		if useList[i].ID == 0 {
			useList[i].ID = i + 1
		}
		useList[i].Random = rand.Float64() * 40
		//m[i].mirrors = &m
	}
	sort.Sort(useList)
}

func ClearUse(id int) {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	//m.inUse = false
	for i := range useList {
		if useList[i].ID == id {
			if useList[i].c_n > 0 {
				log.Println("closing, but found a wait", id)
				// If we have a file waiting for this mirror
				<-useList[i].c
				return
			}
			useList[i].inUse = false
			return
		}
	}
	log.Fatal("Could not find mirror ID", id)
}

func GetMirrorOrQueue(f *FileEntry) (*Mirror, bool) {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	//if len(skip) > 0 {
	//	fmt.Printf("mirror without list: %+v\n", skip)
	//}
	backup := -1
	for i := range useList {
		use := true
		for _, id := range f.skip {
			if id == useList[i].ID {
				use = false
				break
			}
		}
		if use {
			backup = i
		}
		if useList[i].inUse {
			continue
		}
		if use {
			useList[i].inUse = true
			return &useList[i], true
		}
	}

	if backup > -1 {
		useList[backup].queue = append(useList[backup].queue, f)
		return nil, true
	}

	return nil, false
}

func Queue(f *FileEntry) bool {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	//if len(skip) > 0 {
	//	fmt.Printf("mirror without list: %+v\n", skip)
	//}
	for i := range useList {
		use := true
		for _, id := range f.skip {
			if id == useList[i].ID {
				use = false
				break
			}
		}
		if use {
			useList[i].queue = append(useList[i].queue, f)
			return true
		}
	}

	return false
}

func FindStragglers() ([]*FileEntry, *Mirror) {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	for i := range useList {
		if useList[i].inUse {
			continue
		}
		if len(useList[i].queue) > 0 {
			ret := useList[i].queue
			useList[i].queue = []*FileEntry{}
			useList[i].inUse = true
			return ret, &useList[i]
		}
	}
	return nil, nil
}

func GetQueue(id int) []*FileEntry {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	for i := range useList {
		if id == useList[i].ID {
			ret := useList[i].queue
			useList[i].queue = []*FileEntry{}
			return ret
		}
	}
	return []*FileEntry{}
}
