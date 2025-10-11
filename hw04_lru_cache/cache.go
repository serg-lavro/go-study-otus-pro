package hw04lrucache

type Key string

type CacheElement struct {
	ElemValue interface{}
	ElemKey   Key
}

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
		listItem := c.items[key]
		ce := listItem.Value.(CacheElement)
		ce.ElemValue = value
		listItem.Value = ce
		c.queue.MoveToFront(listItem)
		return true
	}

	if c.capacity == c.queue.Len() {
		listItem := c.queue.Back()
		delete(c.items, listItem.Value.(CacheElement).ElemKey)
		c.queue.Remove(listItem)
	}
	ce := CacheElement{ElemValue: value, ElemKey: key}
	listItem := c.queue.PushFront(ce)
	c.items[key] = listItem

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	li, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(li)
		return li.Value.(CacheElement).ElemValue, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
