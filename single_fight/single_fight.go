package single_fight

import "sync"

type call struct {
	val interface{}
	err error
	wg  sync.WaitGroup
}

type Group struct {
	calls map[string]*call
	mu    sync.Mutex
}

func (g *Group) Do(key string, f func() (val interface{}, err error)) (val interface{}, err error) {
	g.mu.Lock()
	if g.calls == nil {
		g.calls = make(map[string]*call)
	}
	if call, ok := g.calls[key]; ok {
		g.mu.Unlock()
		call.wg.Wait()
		return call.val, call.err
	}
	call := &call{}
	call.wg.Add(1)
	g.calls[key] = call
	g.mu.Unlock()

	call.val, call.err = f()
	call.wg.Done()

	g.mu.Lock() //会不会还没取到数据就删除了
	delete(g.calls, key)
	g.mu.Unlock()

	return call.val, call.err
}
