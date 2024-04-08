package lru

import "container/list"

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

//最近最少使用
type LRU struct {
	cache     map[string]*list.Element
	ll        *list.List // 一个双向链表
	maxBytes  int64
	nBytes    int64
	OnEvicted func(key string, value Value) //删除元素时候回调的函数
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *LRU {
	return &LRU{
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		maxBytes:  maxBytes,
		OnEvicted: onEvicted,
	}
}

//删除一个元素
func (l *LRU) RemoveOldest() {
	ele := l.ll.Back()

	if ele != nil {
		l.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(l.cache, kv.key)
		l.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if l.OnEvicted != nil {
			l.OnEvicted(kv.key, kv.value)
		}
	}
	return
}

//新增一个元素
func (l *LRU) Add(key string, val Value) {
	if _, ok := l.cache[key]; ok {
		return
	}
	ele := l.ll.PushFront(&entry{
		key:   key,
		value: val,
	})
	l.cache[key] = ele
	l.nBytes += int64(len(key)) + int64(val.Len())
	for l.maxBytes != 0 && l.maxBytes < l.nBytes {
		l.RemoveOldest()
	}
	return
}

//某个元素移动到队头
func (l *LRU) Get(key string) (value Value, ok bool) {
	if ele, ok := l.cache[key]; ok {
		l.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}
