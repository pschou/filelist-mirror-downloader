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
)

type Mirror struct {
	URL      string
	Latency  float64
	Random   float64
	Failures int
	Client   http.Client
	InUse    bool
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
	fmt.Println(" #  Weight Latency Fails Rand URL")
	for i, e := range m {
		fmt.Printf("%2d) %6.02f %6.02f %6d %6.02f %s\n", i, e.Latency+float64(e.Failures)*20+e.Random, e.Latency, e.Failures, e.Random, e.URL)
	}
}
func (m MirrorList) Shuffle() {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	for i := range m {
		m[i].Random = rand.Float64() * 40
	}
	sort.Sort(m)
}
func (m MirrorList) Pop() *Mirror {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	for i := range m {
		if m[i].InUse == false {
			m[i].InUse = true
			return &m[i]
		}
	}
	return nil
}
