package main

import (
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f GetterFunc = func(key string) ([]byte, error) {
		return []byte(key), nil
	}
	expect := []byte("byte")
	if v, _ := f.Get("byte"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

func TestGet(t *testing.T) {
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}
	loadCounts := make(map[string]int64, len(db))
	group := NewGroup("scores", GetterFunc(
		func(key string) ([]byte, error) {
			loadCounts[key]++
			val, _ := db[key]
			return []byte(val), nil
		},
	), 2<<10)
	for k, v := range db {
		if view, err := group.Get(k); err != nil || view.String() != v { //callback
			t.Fatalf("get callback err")
		}
		if _, err := group.Get(k); err != nil || loadCounts[k] != 1 { //cache
			t.Fatalf("get cacah err")
		}
	}
}
