package cache

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type CacheItem struct {
	Key        string
	Value      []byte
	Headers    map[string][]string
	Expiration time.Time
	LastAccess time.Time
}

type LRUCache struct {
	Capacity  int
	Items     map[string]*list.Element
	EvictList *list.List
	Mutex     sync.RWMutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		Capacity:  capacity,
		Items:     make(map[string]*list.Element),
		EvictList: list.New(),
	}
}

var LRU = NewLRUCache(10)

func (c *LRUCache) Get(key string) (*CacheItem, bool) {

	log.Println("Cache: waiting to aquire lock")
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	log.Println("Cache: lock aquired")

	if elem, ok := c.Items[key]; ok {

		c.EvictList.MoveToFront(elem)

		item := elem.Value.(*CacheItem)
		item.LastAccess = time.Now()

		if time.Now().After(item.Expiration) {
			c.removeElement(elem)
			return nil, false
		}

		return item, true
	}

	return nil, false
}

func (c *LRUCache) Set(key string, value []byte, headers map[string][]string, expiration time.Time) {

	log.Println("Cache: waiting to aquire lock")
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	log.Println("Cache: lock aquired")

	if elem, ok := c.Items[key]; ok {

		c.EvictList.MoveToFront(elem)

		item := elem.Value.(*CacheItem)

		item.Value = value
		item.Headers = headers
		item.Expiration = expiration
		item.LastAccess = time.Now()
	} else {

		if c.EvictList.Len() >= c.Capacity {
			c.removeOldest()
		}

		item := &CacheItem{
			Key:        key,
			Value:      value,
			Headers:    headers,
			Expiration: expiration,
			LastAccess: time.Now(),
		}

		elem := c.EvictList.PushFront(item)
		c.Items[key] = elem
	}
}

func (c *LRUCache) Remove(key string) {

	log.Println("Cache: waiting to aquire lock")
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	log.Println("Cache: lock aquired")

	if elem, found := c.Items[key]; found {
		c.removeElement(elem)
	}
}

func (c *LRUCache) removeOldest() {

	element := c.EvictList.Back()
	if element != nil {
		c.removeElement(element)
	}
}

func (c *LRUCache) removeElement(element *list.Element) {

	c.EvictList.Remove(element)
	item := element.Value.(*CacheItem)

	delete(c.Items, item.Key)
}
