package main

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Response []byte
	ExpireAt time.Time
}

type DNSCache struct {
	mu    sync.RWMutex
	store map[string]CacheEntry
}

var memoryCache = DNSCache{
	store: make(map[string]CacheEntry),
}
