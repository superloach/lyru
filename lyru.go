package lyru

import (
	"container/list"
)

type LRUCache[K comparable, V any] struct {
	cache map[K]*list.Element
	list  *list.List

	capacities [3]float64

	lastValley float64
	lastPeak   float64

	hits   int
	misses int

	CapacityVeloc          float64
	MinCapacityVeloc       float64
	TargetHitRate          float64
	MinCapacity            float64
	MaxCapacity            float64
	CapacityDecel          float64
	EmergencyThreshold     float64
	EmergencyCapacityVeloc float64
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

// NewLRUCache returns a new LRUCache with the given capacity.
func NewLRUCache[K comparable, V any]() *LRUCache[K, V] {
	return &LRUCache[K, V]{
		cache:                  make(map[K]*list.Element),
		list:                   list.New(),
		capacities:             [3]float64{1, 1, 1},
		CapacityVeloc:          0.1,
		lastValley:             0,
		lastPeak:               0,
		TargetHitRate:          0.5,
		hits:                   1,
		misses:                 1,
		MinCapacity:            1,
		MaxCapacity:            0,
		MinCapacityVeloc:       0.001,
		CapacityDecel:          0.9,
		EmergencyThreshold:     0.9,
		EmergencyCapacityVeloc: 0.1,
	}
}

// WithCapacity returns a new LRUCache with the given capacity.
func (c *LRUCache[K, V]) WithCapacity(capacity float64) *LRUCache[K, V] {
	c.capacities[0] = capacity
	c.capacities[1] = capacity
	c.capacities[2] = capacity

	return c
}

// WithCapacityVeloc returns a new LRUCache with the given capacity velocity.

// Put adds a key-value pair to the cache.
func (c *LRUCache[K, V]) Put(key K, value V) {
	// if the key is already in the cache, update the value
	if e, ok := c.cache[key]; ok {
		e.Value.(*entry[K, V]).value = value
		c.list.MoveToFront(e)
	} else {
		// add the new entry to the cache
		e := c.list.PushFront(&entry[K, V]{key, value})
		c.cache[key] = e
	}

	// if the target hit rate is not set, do nothing
	if c.TargetHitRate == 0 {
		return
	}

	justPeaked := false
	justValleyed := false

	if (c.capacities[0] <= c.capacities[1] && c.capacities[1] > c.capacities[2]) ||
		(c.capacities[0] < c.capacities[1] && c.capacities[1] >= c.capacities[2]) {
		// c.capacities[1] is a peak
		c.lastPeak = c.capacities[1]
		justPeaked = true
	} else if c.capacities[0] > c.capacities[1] && c.capacities[1] < c.capacities[2] &&
		c.capacities[0] != c.capacities[1] && c.capacities[1] != c.capacities[2] {
		// c.capacities[1] is a valley
		c.lastValley = c.capacities[1]
		justValleyed = true
	}

	// if the hit rate is above the target hit rate, decrease the capacity
	if c.HitRate() > c.TargetHitRate {
		c.capacities[2] = c.capacities[1]
		c.capacities[1] = c.capacities[0]
		c.capacities[0] = c.capacities[0] - c.CapacityVeloc
	} else if c.HitRate() < c.TargetHitRate {
		c.capacities[2] = c.capacities[1]
		c.capacities[1] = c.capacities[0]
		c.capacities[0] = c.capacities[0] + c.CapacityVeloc
	}

	if c.lastPeak != 0 && c.lastValley != 0 {
		c.capacities[0] = (c.lastPeak + c.lastValley) / 2
		c.capacities[1] = c.capacities[0]
		c.capacities[2] = c.capacities[0]

		c.CapacityVeloc *= c.CapacityDecel

		if justPeaked {
			c.lastValley = 0

		}

		if justValleyed {
			c.lastPeak = 0
		}
	}
	if c.MinCapacity != 0 && c.capacities[0] < c.MinCapacity {
		c.capacities[0] = c.MinCapacity
	}

	if c.MaxCapacity != 0 && c.capacities[0] > c.MaxCapacity {
		c.capacities[0] = c.MaxCapacity
	}

	if c.CapacityVeloc < c.MinCapacityVeloc {
		if c.HitRate() < c.TargetHitRate*c.EmergencyThreshold {
			c.CapacityVeloc = c.EmergencyCapacityVeloc
		} else {
			c.CapacityVeloc = c.MinCapacityVeloc
		}
	}

	// remove the oldest entries until the capacity is reached
	for len(c.cache) >= int(c.capacities[0]) {
		e := c.list.Back()
		c.list.Remove(e)
		delete(c.cache, e.Value.(*entry[K, V]).key)
	}
}

// Get returns the value associated with the given key.
func (c *LRUCache[K, V]) Get(key K) (v V, ok bool) {
	// if the key is not in the cache, return false
	e, ok := c.cache[key]
	if !ok {
		c.misses++
		return v, false
	}

	// increment the hit counter
	c.hits++

	// move the entry to the front of the list
	c.list.MoveToFront(e)
	return e.Value.(*entry[K, V]).value, true
}

// Capacity returns the capacity.
func (c *LRUCache[K, V]) Capacity() float64 {
	return c.capacities[0]
}

// IntCapacity returns the capacity as an int.
func (c *LRUCache[K, V]) IntCapacity() int {
	return int(c.capacities[0])
}

// Hits returns the number of cache hits.
func (c *LRUCache[K, V]) Hits() int {
	return c.hits
}

// Misses returns the number of cache misses.
func (c *LRUCache[K, V]) Misses() int {
	return c.misses
}

// HitRate returns the cache hit rate.
func (c *LRUCache[K, V]) HitRate() float64 {
	return float64(c.hits) / float64(c.hits+c.misses)
}

// LastValley returns the last valley.
func (c *LRUCache[K, V]) LastValley() float64 {
	return c.lastValley
}

// LastPeak returns the last peak.
func (c *LRUCache[K, V]) LastPeak() float64 {
	return c.lastPeak
}
