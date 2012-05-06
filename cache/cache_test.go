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
package cache

import (
	"testing"
	"unsafe"
)

type testCacher struct {
	s string
	i int
}

func (c *testCacher) Size() int {
	return int(unsafe.Sizeof(*c))
}

func TestSimpleCache(t *testing.T) {
	data := []struct {
		k string
		v *testCacher
	}{
		{"k1", &testCacher{"Hello", 1}},
		{"k2", &testCacher{"World", 2}},
	}
	cache := NewCache()

	for _, s := range data {
		cache.Add(s.k, s.v)
	}

	results := []struct {
		k, s string
		i    int
		ok   bool
	}{
		{"k1", "Hello", 1, true},
		{"k2", "World", 2, true},
		{"not here", "", 0, false},
	}

	for _, r := range results {
		c, ok := cache.Get(r.k)
		if ok != r.ok {
			t.Errorf("Get failed for key %q expected success '%t' actual success '%t'", r.k, r.ok, ok)
		}

		if ok {
			v := c.(*testCacher)
			if v.s != r.s || v.i != r.i {
				t.Errorf("Data mismatch, expected (%q, %d) got (%q, %d)", r.s, r.i, v.s, v.i)
			}
		}
	}
}
