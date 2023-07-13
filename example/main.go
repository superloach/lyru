//go:build example

package main

import (
	"log"
	"math/rand"
	"runtime"

	"superloach.xyz/lyru"
)

func main() {
	c := lyru.NewLRUCache[int, struct{}]()

	// add 50 random keys to the cache
	for i := 0; i < 50; i++ {
		// use math/rand to generate random numbers
		key := rand.Intn(1000)

		// add the key to the cache
		c.Put(key, struct{}{})
	}

	// cache random stress test loop
	for i := 0; i < 50000; i++ {
		// use math/rand to generate random numbers
		key := rand.Intn(1000)

		// get the value from the cache
		_, ok := c.Get(key)

		// if the key was not found in the cache, add it
		if !ok {
			c.Put(key, struct{}{})

			m := runtime.MemStats{}
			runtime.ReadMemStats(&m)

			runtime.GC()

			m2 := runtime.MemStats{}
			runtime.ReadMemStats(&m2)

			log.Printf(
				"[%5d] hits: %06d, misses: %04d, hit rate: %f, capacity: %f, capacityVeloc: %f, alloc: %d, alloc2: %d, last valley: %f, last peak: %f",
				i, c.Hits(), c.Misses(), c.HitRate(), c.Capacity(), c.CapacityVeloc, m.Alloc, m2.Alloc, c.LastValley(), c.LastPeak(),
			)
		}
	}

	// start stressing with misses
	for i := 0; i < 50000; i++ {
		// use math/rand to generate random numbers
		key := rand.Intn(1000) + 1000

		// get the value from the cache
		_, ok := c.Get(key)

		// if the key was not found in the cache, add it
		if !ok {
			c.Put(key, struct{}{})

			m := runtime.MemStats{}
			runtime.ReadMemStats(&m)

			runtime.GC()

			m2 := runtime.MemStats{}
			runtime.ReadMemStats(&m2)

			log.Printf(
				"[%5d] hits: %06d, misses: %04d, hit rate: %f, capacity: %f, capacityVeloc: %f, alloc: %d, alloc2: %d, last valley: %f, last peak: %f",
				i, c.Hits(), c.Misses(), c.HitRate(), c.Capacity(), c.CapacityVeloc, m.Alloc, m2.Alloc, c.LastValley(), c.LastPeak(),
			)
		}
	}
}
