package main

import (
	"log"
	"sync"

	"cache.example/single_fight"
)

// 定义一个函数类型 F，并且实现接口 A 的方法，然后在这个方法中调用自己。
// 这是 Go 语言中将其他函数（参数返回值定义与 F 一致）转换为接口 A 的常用技巧。
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu       sync.Mutex
	groupMap map[string]*Group
)

type Group struct {
	name      string
	mainCache cache
	getter    Getter
	peers     PeerPicker
	onceCall  *single_fight.Group
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

//注册一致性哈希
// func r

func NewGroup(name string,
	getter Getter,
	cacheBytes int64) *Group {
	if getter == nil {
		panic("getter is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
		onceCall: &single_fight.Group{},
	}
	if groupMap == nil {
		groupMap = make(map[string]*Group)
	}
	groupMap[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	return groupMap[name]
}

func (g *Group) Get(key string) (val ByteView, err error) {
	b, ok := g.mainCache.get(key)
	if ok {
		return b, nil
	}
	//回源
	return g.load(key)
}

func (g *Group) Add(key string, val ByteView) {
	g.mainCache.Add(key, val)
	return
}

func (g *Group) load(key string) (value ByteView, err error) {

	viewi, err := g.onceCall.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

//调用回调函数
func (g *Group) getLocally(key string) (ByteView, error) {
	val, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	byteView := ByteView{
		b: val,
	}
	g.Add(key, byteView)
	return byteView, nil
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	val, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	byteView := ByteView{
		b: val,
	}
	return byteView, nil
}
