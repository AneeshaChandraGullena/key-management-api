// Package cache
// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.
package cache

import (
	"crypto/sha256"
	"sync"
	"time"
)

type cacheEntry struct {
	used       bool   // Has this entry been recently used?
	token      string //The token
	role       string //The role of the user
	expiration int64  //Expiration of the token as a UNIX time stamp.
}

type cache struct {
	entries   []cacheEntry   //The actual cache entries
	hash      map[string]int //Reverse lookup table to speed up get times.
	length    int            //Size of our cache
	clockHand int            //Clock hand used for approximating LRU.
	mutex     *sync.RWMutex
}

const (
	cacheEntrySize         = 1 << 11  // 2 kb for a JWT token.
	availableSystemMemory  = 32882120 // 32 gb in kb. Taken from /proc/meminfo
	percentOfMemoryToCache = .01      // Percent of the memory to devote to the token cache
)

var cacheInstance *cache

func initCache(size int) {
	cacheInstance = new(cache)
	cacheInstance.mutex = new(sync.RWMutex)
	cacheInstance.entries = make([]cacheEntry, size)
	cacheInstance.length = size
	cacheInstance.hash = make(map[string]int)
}

func InitCacheCompute() {
	cacheInstance = new(cache)
	cacheInstance.mutex = new(sync.RWMutex)
	size := computeSize()
	cacheInstance.entries = make([]cacheEntry, size)
	cacheInstance.hash = make(map[string]int)
	cacheInstance.length = size
}

//Computes the size of cache based off the resources available.
func computeSize() int {
	availableCacheMemory := availableSystemMemory * percentOfMemoryToCache
	return int(availableCacheMemory / cacheEntrySize)
}

//Insert a token and it's expiration date into the cache
func Insert(token string, expiration int64, defaultExpr bool, role string) bool {
	exists, _ := Get(token)
	if exists {
		return false
	}
	cacheInstance.mutex.Lock()
	defer cacheInstance.mutex.Unlock()
	var inUse bool
	now := time.Now().Unix()
	if defaultExpr {
		expiration = time.Now().Add(time.Minute * 10).Unix()
	}
	encryptedToken := encrypt(token)
	for i := 0; i <= cacheInstance.length; i++ {
		inUse = cacheInstance.entries[cacheInstance.clockHand].used
		expir := cacheInstance.entries[cacheInstance.clockHand].expiration
		if inUse == false || expir <= now {
			//Found an empty cache spot
			//Remove the value from the reverse lookup table.
			delete(cacheInstance.hash, cacheInstance.entries[cacheInstance.clockHand].token)
			//Insert the new value into the cache
			cacheInstance.entries[cacheInstance.clockHand].used = true
			cacheInstance.entries[cacheInstance.clockHand].token = encryptedToken
			cacheInstance.entries[cacheInstance.clockHand].expiration = expiration
			cacheInstance.entries[cacheInstance.clockHand].role = role
			cacheInstance.hash[encryptedToken] = cacheInstance.clockHand
			cacheInstance.clockHand = (cacheInstance.clockHand + 1) % cacheInstance.length
			return true
		}
		//Not an empty cache spot. Mark as unused
		cacheInstance.entries[cacheInstance.clockHand].used = false
		cacheInstance.clockHand = (cacheInstance.clockHand + 1) % cacheInstance.length
	}
	//Never reaches here.
	return false
}

//Check if a token is in the cache.
func Get(token string) (bool, string) {
	cacheInstance.mutex.RLock()
	defer cacheInstance.mutex.RUnlock()
	now := time.Now().Unix()
	encryptedToken := encrypt(token)
	if i, ok := cacheInstance.hash[encryptedToken]; ok {
		if cacheInstance.entries[i].expiration > now {
			cacheInstance.entries[i].used = true
			return true, cacheInstance.entries[i].role
		}
	}
	return false, ""
}

func encrypt(token string) string {
	bytes := sha256.Sum256([]byte(token))
	return string(bytes[:])
}
