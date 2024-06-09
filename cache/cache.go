package cache

import (
	"sync"
	"time"
)

var Instance *Cache

type CacheItem struct {
	Value           float64
	LastAccessDTSec int64
	TimeToLiveSec   int64
}

type Cache struct {
	mtx            sync.Mutex
	items          map[string]*CacheItem
	lastClearDTStr string
	clearCount     int
	getCount       int
	setCount       int
	removeCount    int
}

type CacheState struct {
	CountOfItems   int
	LastClearDTStr string
	ClearCount     int
	GetCount       int
	SetCount       int
	RemoveCount    int
}

func init() {
	Instance = NewCache()
}

func NewCache() *Cache {
	var c Cache
	c.items = make(map[string]*CacheItem)
	return &c
}

func (c *Cache) Start() {
	go c.thClear()
}

func (c *Cache) GetState() *CacheState {
	var state CacheState
	c.mtx.Lock()
	state.CountOfItems = len(c.items)
	state.ClearCount = c.clearCount
	state.GetCount = c.getCount
	state.SetCount = c.setCount
	state.RemoveCount = c.removeCount
	c.mtx.Unlock()
	return &state
}

func (c *Cache) thClear() {
	for {
		time.Sleep(10 * time.Second)
		c.clearCount++
		c.lastClearDTStr = time.Now().UTC().Format("2006-01-02 15:04:05")
		itemsToRemove := make([]string, 0)
		nowSec := time.Now().Unix()

		c.mtx.Lock()
		for key, value := range c.items {
			age := nowSec - value.LastAccessDTSec
			if age > value.TimeToLiveSec {
				itemsToRemove = append(itemsToRemove, key)
			}
		}
		for _, key := range itemsToRemove {
			delete(c.items, key)
			c.removeCount++
		}
		c.mtx.Unlock()
	}
}

func (c *Cache) Get(id string) *CacheItem {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.getCount++
	if item, ok := c.items[id]; ok {
		item.LastAccessDTSec = time.Now().Unix()
		return item
	}
	return nil
}

func (c *Cache) Set(id string, value float64) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.setCount++
	var item CacheItem
	item.Value = value
	item.LastAccessDTSec = time.Now().Unix()
	item.TimeToLiveSec = 5 * 60
	c.items[id] = &item
}
