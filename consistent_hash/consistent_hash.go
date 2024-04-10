package consistent_hash

import (
	"hash/crc32"
	"log"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int //Sorted
	hashMap  map[int]string
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys []string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			k := strconv.Itoa(i) + key
			hash := int(m.hash([]byte(k)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	log.Printf("m.hashMap:%v", m.hashMap)
	sort.Ints(m.keys)
	return
}

func (m *Map) Get(val string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := m.hash([]byte(val))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= int(hash)
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
