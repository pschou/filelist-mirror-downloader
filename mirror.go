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
	inUse    bool
	//mirrors  *MirrorList
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
	fmt.Println(" #  Weight Latency Fails   Rand InUse URL")
	for _, e := range m {
		fmt.Printf("%2d) %6.02f %6.02f %6d %6.02f %t %s\n", e.ID, e.Latency+float64(e.Failures)*20+e.Random, e.Latency, e.Failures, e.Random, e.inUse, e.URL)
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
			useList[i].inUse = false
			return
		}
	}
	log.Fatal("Could not find mirror ID", id)
}
func PopWithout(skip []int) *Mirror {
	MirrorListSync.Lock()
	defer MirrorListSync.Unlock()
	//fmt.Println("mirror list: %+v\n", m)
	for i := range useList {
		use := true
		for id := range skip {
			if id != useList[i].ID && useList[i].inUse == false {
				use = false
				break
			}
		}
		if use {
			useList[i].inUse = true
			return &useList[i]
		}
	}
	return nil
}
