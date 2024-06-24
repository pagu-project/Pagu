package cache

import (
	"testing"
	"time"
)

var keyVal = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
}

func TestCache(t *testing.T) {
	cache := NewBasic[string, int](0)

	t.Run("add method", func(t *testing.T) {
		for k, v := range keyVal {
			if ok := cache.Add(k, v, 0); !ok {
				t.Errorf("failed add %s to cache", k)
			}
		}
	})

	t.Run("get method", func(t *testing.T) {
		for k := range keyVal {
			val, ok := cache.Get(k)
			if !ok {
				t.Errorf("Get method failed. Expected Value: 1, Got: %v, Ok: %v", val, ok)
			}
		}
	})

	t.Run("update", func(t *testing.T) {
		for k := range keyVal {
			ok := cache.Update(k, 5, 3*time.Second)
			if !ok {
				t.Errorf("failed update %s to cache", k)
			}
			val, ok := cache.Get(k)
			if !ok {
				t.Errorf("Get method failed. Expected Value: 1, Got: %v, Ok: %v", val, ok)
			}

			break
		}
	})

	t.Run("list keys", func(t *testing.T) {
		if len(cache.Keys()) != len(keyVal) {
			t.Error("list keys is invalid")
		}
	})

	t.Run("exists keys", func(t *testing.T) {
		for k := range keyVal {
			if ok := cache.Exists(k); !ok {
				t.Errorf("key %s not exists", k)
			}
		}
	})

	t.Run("delete items", func(t *testing.T) {
		for k := range keyVal {
			if ok := cache.Delete(k); !ok {
				t.Errorf("failed to delete %s", k)
			}
		}
	})
}

func TestCacheWithExpiry(t *testing.T) {
	cache := NewBasic[string, int](1 * time.Second)

	t.Run("add cache for expiration", func(t *testing.T) {
		step := 5
		for k, v := range keyVal {
			if ok := cache.Add(k, v, time.Duration(step)*time.Second); !ok {
				t.Errorf("failed add %s to cache", k)
			}
			step += 2
		}

		// Wait for cache entries to expire
		time.Sleep(15 * time.Second)

		// Check that expired entries are no longer in the cache
		for k := range keyVal {
			if cache.Exists(k) {
				t.Errorf("expected %s to be expired and removed from cache", k)
			}
		}
	})
}

func BenchmarkAdd(b *testing.B) {
	b.ReportAllocs()
	cache := NewBasic[int, int](0) // or NewBasic[int, int](0, WithCompressor())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Add(i, i, 0)
	}
}

func BenchmarkGet(b *testing.B) {
	b.ReportAllocs()
	cache := NewBasic[int, int](0) // or NewBasic[int, int](0, WithCompressor())

	// Pre-populate the cache
	for i := 0; i < 1000; i++ {
		cache.Add(i, i, 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Get(i)
	}
}

func BenchmarkUpdate(b *testing.B) {
	b.ReportAllocs()
	cache := NewBasic[int, int](0) // or NewBasic[int, int](0, WithCompressor())

	// Pre-populate the cache
	for i := 0; i < 1000; i++ {
		cache.Add(i, i, 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Update(i, i+1, 0)
	}
}
