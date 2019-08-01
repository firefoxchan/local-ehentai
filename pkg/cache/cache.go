package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

type Catch struct {
	data map[string]*cNode
	mu   sync.RWMutex
}

type cNode struct {
	value    interface{}
	expireAt int64
}

func NewCache(cleanInterval time.Duration) *Catch {
	c := &Catch{
		data: map[string]*cNode{},
		mu:   sync.RWMutex{},
	}
	go c.run(cleanInterval)
	return c
}

func (c *Catch) run(cleanInterval time.Duration) {
	for {
		select {
		case <-time.After(cleanInterval):
		}
		now := time.Now()
		c.mu.Lock()
		for k, v := range c.data {
			if time.Unix(atomic.LoadInt64(&v.expireAt), 0).Before(now) {
				delete(c.data, k)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Catch) Set(key string, value interface{}, duration time.Duration) {
	cn := &cNode{
		value:    value,
		expireAt: time.Now().Add(duration).Unix(),
	}
	c.mu.Lock()
	c.data[key] = cn
	c.mu.Unlock()
}

func (c *Catch) Get(key string, renew time.Duration) (interface{}, bool) {
	c.mu.RLock()
	cn, ok := c.data[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	now := time.Now()
	if time.Unix(cn.expireAt, 0).Before(now) {
		return nil, false
	}
	if renew > 0 {
		atomic.StoreInt64(&cn.expireAt, now.Add(renew).Unix())
	}
	return cn.value, true
}
