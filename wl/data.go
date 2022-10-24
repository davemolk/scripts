package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

type WordMap struct {
	mu    sync.Mutex
	words map[string]int
}

func newWordMap() *WordMap {
	return &WordMap{
		words: make(map[string]int),
	}
}

func (wm *WordMap) add(w string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.words[w]++
}

func (wm *WordMap) delete(key string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	delete(wm.words, key)
}

func (wm *WordMap) sort() []string {
	keys := make([]string, 0, len(wm.words))
	for key := range wm.words {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return wm.words[keys[i]] > wm.words[keys[j]]
	})

	return keys
}

func (wm *WordMap) write(keys []string, name string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, key := range keys {
		fmt.Fprintf(f, "%s: %d\n", key, wm.words[key])
	}
}
