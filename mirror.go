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
	ID       int
	URL      string
	Latency  float64
	Random   float64
	Failures int
	Client   http.Client
	InUse    bool
	mirrors  MirrorList
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
	fmt.Println(" #  Weight Latency Fails   Rand Use  URL")
	for i, e := range m {
		fmt.Printf("%2d) %6.02f %6.02f %6d %6.02f %t %s\n", i, e.Latency+float64(e.Failures)*20+e.Random, e.Latency, e.Failures, e.Random, e.InUse, e.URL)
	}
}
func (m MirrorList) Shuffle() {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	for i := range m {
		if m[i].ID == 0 {
			m[i].ID = i + 1
		}
		m[i].Random = rand.Float64() * 40
		m[i].mirrors = m
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
func (m Mirror) ClearUse(id int) {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	for i := range m.mirrors {
		if m.mirrors[i].ID == id {
			m.mirrors[i].InUse = false
			return
		}
	}
}
func (m MirrorList) PopWithout(skip []int) *Mirror {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	//fmt.Println("mirror list: %+v\n", m)
poploop:
	for i := range m {
		for id := range skip {
			if id == m[i].ID {
				continue poploop
			}
		}
		if m[i].InUse == false {
			m[i].InUse = true
			return &m[i]
		}
	}
	return nil
}
