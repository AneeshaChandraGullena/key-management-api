// Package cache
// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.
package cache

import "testing"
import "strconv"

func TestSizeCalculation(t *testing.T) {
	calculatedSize := computeSize()
	availableCacheMemory := availableSystemMemory * percentOfMemoryToCache
	testSize := int(availableCacheMemory / cacheEntrySize)
	if testSize != calculatedSize {
		t.Errorf("Expected size %d got size %d\n", testSize, calculatedSize)
		t.Fail()
	}
}

func TestInit(t *testing.T) {
	initCache(2)
}

func TestInitCompute(t *testing.T) {
	InitCacheCompute()
}

func TestCacheInsert(t *testing.T) {
	entered := Insert("12345678", 2573783537, false, "")
	if !entered {
		t.Error("Cache entry couldn't be entered")
		t.Fail()
	}
}

func TestCacheInsertDuplicate(t *testing.T) {
	entered := Insert("1234567890", 2573783537, false, "")
	if !entered {
		t.Error("Cache entry couldn't be entered")
		t.Fail()
	}
	entered = Insert("1234567890", 2573783537, false, "")
	if entered {
		t.Error("Cache allowed a duplicate token.")
		t.Fail()
	}
}

func TestInsertandGet(t *testing.T) {
	entered := Insert("123", 2573783537, false, "")
	if !entered {
		t.Error("Cache entry couldn't be entered")
		t.Fail()
	}
	exists, _ := Get("123")
	if !exists {
		t.Error("Cache entry couldn't be found")
		t.Fail()
	}
}

func TestCacheGetNoSuchToken(t *testing.T) {
	exists, _ := Get("abc1234")
	if exists {
		t.Error("Cache said it has token but shouldn't")
		t.Fail()
	}
}

func TestCacheGetExpired(t *testing.T) {
	entered := Insert("1234", 0, false, "")
	if !entered {
		t.Error("Cache entry couldn't be entered")
		t.Fail()
	}
	exists, _ := Get("1234")
	if exists {
		t.Error("Cache said it has token but it's expired")
		t.Fail()
	}
}

func TestExpiredTokenReplacement(t *testing.T) {
	entered := Insert("zxy", 0, false, "")
	if !entered {
		t.Error("Cache entry couldn't be entered")
		t.Fail()
	}

	entered = Insert("xyz", 2573783537, false, "")
	if !entered {
		t.Error("Cache entry couldn't be entered")
		t.Fail()
	}

	entered = Insert("xyz2", 2573783537, false, "")
	if !entered {
		t.Error("Cache entry couldn't be entered")
		t.Fail()
	}

	exists, _ := Get("zxy")
	if exists {
		t.Error("Cache said it has token but it's expired")
		t.Fail()
	}

	exists, _ = Get("xyz")
	if !exists {
		t.Error("Cache lost valid token")
		t.Fail()
	}

	exists, _ = Get("xyz2")
	if !exists {
		t.Error("Cache lost valid token")
		t.Fail()
	}
}

func TestCacheGetDefaultExpr(t *testing.T) {
	entered := Insert("mytesttoken", 0, true, "")
	if !entered {
		t.Error("Cache entry couldn't be entered")
		t.Fail()
	}
	exists, _ := Get("mytesttoken")
	if !exists {
		t.Error("Couldn't find a token with a default expiration")
		t.Fail()
	}
}

func TestCacheGetRoleTest(t *testing.T) {
	role := "testRole"
	entered := Insert("mytesttoken22", 0, true, role)
	if !entered {
		t.Error("Cache entry couldn't be entered")
		t.Fail()
	}
	_, returnedRole := Get("mytesttoken22")
	if returnedRole != role {
		t.Error("Role doesn't match expected role.")
		t.Fail()
	}
}

func BenchmarkCache(b *testing.B) {
	size := 1000
	initCache(size)
	tokenArray := make([]string, size)
	for i := 0; i < size; i++ {
		tokenArray[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < size; j++ {
			Insert(tokenArray[j], 0, true, "")
		}
		for j := 0; j < size; j++ {
			Get(tokenArray[j])
		}
	}
}
