package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if i, ok := l.items[key]; ok {
		l.queue.MoveToFront(i)
		i.Value.(*cacheItem).value = value
		return true
	}

	l.items[key] = l.queue.PushFront(&cacheItem{key, value})
	if l.capacity < l.queue.Len() {
		l.Clear()
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if i, ok := l.items[key]; ok {
		l.queue.MoveToFront(i)
		return i.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	if i := l.queue.Back(); i != nil {
		l.queue.Remove(i)
		delete(l.items, i.Value.(*cacheItem).key)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
