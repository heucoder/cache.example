package consistent_hash

import (
	"strconv"
	"testing"
)

func TestConsistentHash(t *testing.T) {
	var keyNodeMap = map[string]string{
		"1": "2",
		"2": "2",
		"3": "4",
		"4": "4",
		"5": "2",
	}
	var h Hash = func(data []byte) uint32 {
		idx, _ := strconv.Atoi(string(data))
		return uint32(idx)
	}
	m := New(2, h)
	m.Add([]string{"4", "2"}...) //4 14 2 12
	for k, v := range keyNodeMap {
		node := m.Get(k)
		if node != v {
			t.Fatalf("get node k:%s v:%s node:%s", k, v, node)
		}
	}
}
