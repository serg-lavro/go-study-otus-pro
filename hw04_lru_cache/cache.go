package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	_, ok := c.items[key]
	if ok {
		cacheElem := c.items[key]
		cacheElem.Value = value
		c.queue.MoveToFront(cacheElem)
		return true
	}

	if c.capacity == c.queue.Len() {
		cacheElem := c.queue.Back()
		c.queue.Remove(cacheElem)
		delete(c.items, cacheElem.ItemKey)
	}
	cacheElem := c.queue.PushFront(value)
	cacheElem.ItemKey = key
	c.items[key] = cacheElem

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	elem, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(elem)
		return elem.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
