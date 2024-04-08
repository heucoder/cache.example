package lru

import (
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(0, nil)
	lru.Add("1", String("123456789"))
	v, ok := lru.Get("1")
	if !ok {
		t.Fatalf("key not found")
	}
	s := v.(String)
	if s != "123456789" {
		t.Fatalf("val not equal")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"
	cap := len(v1) + len(v2) + len(k1) + len(k2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	if len(lru.cache) != 2 {
		t.Fatalf("未能自动移除")
	}
	v, ok := lru.Get("k1")
	if ok {
		t.Fatalf("移除k不对")
	}
	if v != nil {
		t.Fatalf("移除v不对")
	}

}

func TestOnEvicted(t *testing.T) {
	deletekeyList := []string{}
	callBack := func(key string, value Value) {
		deletekeyList = append(deletekeyList, key)
	}
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"
	lru := New(1, callBack)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	if len(deletekeyList) != 3 {
		t.Fatalf("过期移除未生效")
	}
}
