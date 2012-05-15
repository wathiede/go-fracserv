// Copyright 2012 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package fracserv

import (
	"encoding/json"
	"encoding/gob"
	"log"
	"flag"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	bookmarkFn = flag.String("bookmarkFn", "/tmp/bookmark.db",
		"location to store bookmarked links")
)

func init() {
	b := NewBookmarks()
	http.HandleFunc("/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		b.ListHandler(w, r)
	})
	http.HandleFunc("/bookmarks/add", func(w http.ResponseWriter, r *http.Request) {
		b.AddHandler(w, r)
	})
}

type Bookmark struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Added time.Time `json:"added"`
}

type Bookmarks struct {
	Bookmarks []Bookmark
	addCh chan Bookmark
	mu   sync.RWMutex
}

func NewBookmarks() *Bookmarks {
	b := &Bookmarks{
		Bookmarks: make([]Bookmark, 0),
		addCh: make(chan Bookmark, 1),
	}
	go b.Loop()
	return b
}

func (b *Bookmarks) Add(bookmark Bookmark) {
	log.Print("Adding bookmark ", bookmark)
	b.addCh <- bookmark
}

func (b *Bookmarks) Load(fn string) error {
	log.Print("Loading bookmarks from ", fn)

	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	gf := gob.NewDecoder(f)
	err = gf.Decode(&b.Bookmarks)
	return err
}

func (b *Bookmarks) Save(fn string) error {
	log.Print("Saving bookmarks to ", fn)
	f, err := os.OpenFile(fn, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	gf := gob.NewEncoder(f)
	err = gf.Encode(b.Bookmarks)
	return err
}

func (b *Bookmarks) Loop() {
	// TODO(wathiede) if this is called from init, it won't get the path
	// specified by flag.bookmarkFn
	err := b.Load(*bookmarkFn)
	if err != nil {
		log.Print("Error loading bookmarks ", err)
	}
	for {
		select {
		case bookmark := <-b.addCh:
			b.mu.Lock()
			b.Bookmarks = append(b.Bookmarks, bookmark)
			// TODO This doesn't really scale
			err = b.Save(*bookmarkFn)
			if err != nil {
				log.Print("Error saving bookmarks ", err)
			}
			b.mu.Unlock()
		}
	}
}

func (b *Bookmarks) AddHandler(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	name := q.Get("name")
	url := q.Get("url")
	if url == "" || name == "" {
		http.Error(w, "missing url or name: " + req.URL.String(),
			http.StatusInternalServerError)
		return
	}
	b.Add(Bookmark{
		Name: q.Get("name"),
		Url: q.Get("url"),
		Added: time.Now(),
	})

	http.Redirect(w, req, "/", http.StatusMovedPermanently)
}

func (b *Bookmarks) ListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b.mu.Lock()
	defer b.mu.Unlock()
	j, err := json.Marshal(b.Bookmarks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(j)
}
