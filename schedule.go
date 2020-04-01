package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type CouterMap interface {
	CompareOrIncrement(string, int) bool
}

type GlobalLockCounterMap struct {
	mu      sync.Mutex
	Counter map[string]int
}

func (m *GlobalLockCounterMap) CompareOrIncrement(key string, threshold int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	count, ok := m.Counter[key]
	if ok {
		if count < threshold {
			count += 1
			return true
		} else {
			return false
		}
	} else {
		m.Counter[key] = 1
		return true
	}
}
//type HostCounter struct {
//	Count int
//	mu    sync.Mutex
//}
//
//type RateLimiter struct {
//	Counter map[string]*HostCounter
//}
//
//func (r *RateLimiter) Increment(hostname string, threshold int) bool {
//	counter, ok := r.Counter[hostname]
//	if !ok {
//		r.Counter[hostname].Count = 1
//
//	} else {
//		counter.mu.Lock()
//		counter.Count += 1
//		counter.mu.Unlock()
//	}
//	atomic.AddUint64(&ops, 1)
//}

// StartSpider starts the spider. Call this function in a goroutine.
func Limit() {
	nextTime := time.Now().Truncate(time.Second)
	for {
		nextTime = nextTime.Add(time.Second)
		time.Sleep(time.Until(nextTime))
		go fmt.Printf("Spider : %s\n", time.Now())
	}
}

func main() {
	RateLimiter()
}
