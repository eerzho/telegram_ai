package lru

import "sync"

type node[T any] struct {
	key  string
	val  T
	prev *node[T]
	next *node[T]
}

type Cache[T any] struct {
	mu       sync.Mutex
	capacity int
	data     map[string]*node[T]
	head     *node[T]
	tail     *node[T]
}

func NewCache[T any](capacity int) *Cache[T] {
	if capacity < 0 {
		capacity = 1
	}

	head := &node[T]{}
	tail := &node[T]{}

	head.next = tail
	tail.prev = head

	return &Cache[T]{
		capacity: capacity,
		data:     make(map[string]*node[T]),
		head:     head,
		tail:     tail,
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero T
	n, ok := c.data[key]
	if !ok {
		return zero, false
	}
	c.remove(n)
	c.addToHead(n)
	return n.val, true
}

func (c *Cache[T]) Put(key string, val T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if n, ok := c.data[key]; ok {
		c.remove(n)
		c.addToHead(n)
		n.val = val
		return
	}

	n := &node[T]{key: key, val: val}

	if len(c.data) >= c.capacity {
		delete(c.data, c.tail.prev.key)
		c.remove(c.tail.prev)
	}

	c.data[key] = n
	c.addToHead(n)
}

func (c *Cache[T]) remove(n *node[T]) {
	prev := n.prev
	next := n.next

	prev.next = next
	next.prev = prev
}

func (c *Cache[T]) addToHead(n *node[T]) {
	next := c.head.next

	c.head.next = n
	n.prev = c.head

	n.next = next
	next.prev = n
}
