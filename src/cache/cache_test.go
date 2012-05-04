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
		i int
		ok bool
	}{
		{"k1", "Hello", 1, true},
		{"k2", "World", 2, true},
		{"not here", "", 0, false},
	}

	for _, r := range results {
		c, ok := cache.Get(r.k);
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
