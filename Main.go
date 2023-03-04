package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// defines the interface that any mutex lock must implement.
type Mutex interface {
	Lock()
	Unlock()
}

// a pointer to Locker implements mutex interface.
type Locker struct {
	// toggles b/w only 0 and 1.
	// 0 means lock is available.
	// 1 means lock is acquired.
	lock uintptr
}

func (l *Locker) Lock() {
	for {
		if atomic.CompareAndSwapUintptr(&l.lock, 0, 1) {
			break
		}
	}
}

func (l *Locker) Unlock() {
	atomic.StoreUintptr(&l.lock, 0)
}

type Counter struct {
	i  int
	mu Mutex
}

func (c *Counter) incr(wg *sync.WaitGroup, safe bool) {
	defer wg.Done()

	if safe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}
	c.i++
}

func main() {
	counter := Counter{
		i: 0,
		mu: &Locker{
			lock: 0,
		},
	}

	increments := 1000

	var wg sync.WaitGroup
	wg.Add(increments)
	for i := 0; i < increments; i++ {
		// if second argument is false, value of counter.i is
		// not guaranteed to be the same as increments.
		go counter.incr(&wg, true)
	}
	wg.Wait()
	fmt.Println(counter.i)
}
