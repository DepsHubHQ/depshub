package sources

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNewFileCache(t *testing.T) {
	tests := []struct {
		name      string
		cacheName string
		wantErr   bool
	}{
		{
			name:      "valid cache name",
			cacheName: "test-cache",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := NewFileCache(tt.cacheName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if cache == nil {
				t.Error("NewFileCache() returned nil cache")
			}

			// Cleanup
			if cache != nil {
				os.Remove(cache.filename)
			}
		})
	}
}

func TestFileCache_SetGet(t *testing.T) {
	cache, err := NewFileCache("test-cache")
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer os.Remove(cache.filename)

	tests := []struct {
		name        string
		key         string
		value       interface{}
		expiration  time.Duration
		wantErr     bool
		wantExists  bool
		checkExpiry bool
	}{
		{
			name:       "string value",
			key:        "string-key",
			value:      "test-value",
			expiration: time.Hour,
			wantExists: true,
		},
		{
			name: "map value",
			key:  "map-key",
			value: map[string]string{
				"key": "value",
			},
			expiration: time.Hour,
			wantExists: true,
		},
		{
			name:        "expired item",
			key:         "expired-key",
			value:       "expired-value",
			expiration:  -time.Hour, // Already expired
			wantExists:  false,
			checkExpiry: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set value
			err := cache.Set(tt.key, tt.value, tt.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileCache.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Get value
			var got interface{}
			exists, err := cache.Get(tt.key, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileCache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if exists != tt.wantExists {
				t.Errorf("FileCache.Get() exists = %v, want %v", exists, tt.wantExists)
				return
			}

			if exists && !tt.checkExpiry {
				// Check if retrieved value matches stored value
				switch v := tt.value.(type) {
				case string:
					if got != v {
						t.Errorf("FileCache.Get() = %v, want %v", got, v)
					}
				case map[string]string:
					gotMap, ok := got.(map[string]interface{})
					if !ok {
						t.Errorf("Failed to convert got to map[string]interface{}")
						return
					}
					for k, want := range v {
						if got := gotMap[k]; got != want {
							t.Errorf("FileCache.Get() map[%s] = %v, want %v", k, got, want)
						}
					}
				}
			}
		})
	}
}

func TestFileCache_Delete(t *testing.T) {
	cache, err := NewFileCache("test-cache")
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer os.Remove(cache.filename)

	// Set a value
	key := "test-key"
	value := "test-value"
	if err := cache.Set(key, value, time.Hour); err != nil {
		t.Fatalf("Failed to set cache value: %v", err)
	}

	// Delete the value
	if err := cache.Delete(key); err != nil {
		t.Errorf("FileCache.Delete() error = %v", err)
		return
	}

	// Verify it's gone
	var got string
	exists, err := cache.Get(key, &got)
	if err != nil {
		t.Errorf("FileCache.Get() error = %v", err)
		return
	}
	if exists {
		t.Error("FileCache.Get() returned exists = true after Delete")
	}
}

func TestFileCache_Clear(t *testing.T) {
	cache, err := NewFileCache("test-cache")
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer os.Remove(cache.filename)

	// Set multiple values
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	for k, v := range testData {
		if err := cache.Set(k, v, time.Hour); err != nil {
			t.Fatalf("Failed to set cache value: %v", err)
		}
	}

	// Clear the cache
	if err := cache.Clear(); err != nil {
		t.Errorf("FileCache.Clear() error = %v", err)
		return
	}

	// Verify all values are gone
	for k := range testData {
		var got string
		exists, err := cache.Get(k, &got)
		if err != nil {
			t.Errorf("FileCache.Get() error = %v", err)
			return
		}
		if exists {
			t.Errorf("FileCache.Get() returned exists = true after Clear for key %s", k)
		}
	}
}

func TestFileCache_Persistence(t *testing.T) {
	cacheName := "test-persistence-cache"
	cache, err := NewFileCache(cacheName)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer os.Remove(cache.filename)

	// Set a value
	key := "test-key"
	value := "test-value"
	if err := cache.Set(key, value, time.Hour); err != nil {
		t.Fatalf("Failed to set cache value: %v", err)
	}

	// Create a new cache instance with the same name
	cache2, err := NewFileCache(cacheName)
	if err != nil {
		t.Fatalf("Failed to create second cache: %v", err)
	}

	// Verify the value exists in the new instance
	var got string
	exists, err := cache2.Get(key, &got)
	if err != nil {
		t.Errorf("FileCache.Get() error = %v", err)
		return
	}
	if !exists {
		t.Error("FileCache.Get() returned exists = false for persisted value")
		return
	}
	if got != value {
		t.Errorf("FileCache.Get() = %v, want %v", got, value)
	}
}

func TestFileCache_Concurrent(t *testing.T) {
	cache, err := NewFileCache("test-concurrent-cache")
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer os.Remove(cache.filename)

	done := make(chan bool)
	const goroutines = 10

	// Concurrent writes
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			key := fmt.Sprintf("key-%d", id)
			value := fmt.Sprintf("value-%d", id)
			err := cache.Set(key, value, time.Hour)
			if err != nil {
				t.Errorf("Concurrent Set failed: %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < goroutines; i++ {
		<-done
	}

	// Verify all values were written correctly
	for i := 0; i < goroutines; i++ {
		key := fmt.Sprintf("key-%d", i)
		expectedValue := fmt.Sprintf("value-%d", i)
		var got string
		exists, err := cache.Get(key, &got)
		if err != nil {
			t.Errorf("Get failed for key %s: %v", key, err)
			continue
		}
		if !exists {
			t.Errorf("Value not found for key %s", key)
			continue
		}
		if got != expectedValue {
			t.Errorf("Got %s, want %s for key %s", got, expectedValue, key)
		}
	}
}
