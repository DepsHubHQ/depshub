// The cache file will be stored in:

// Windows: %USERPROFILE%\.cache\depshub\
// Linux/macOS: ~/.cache/depshub/

package sources

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CacheItem struct {
	Value      []byte    `json:"value"` // Store as raw JSON bytes
	ExpiresAt  time.Time `json:"expires_at"`
	CreateTime time.Time `json:"create_time"`
}

type FileCache struct {
	filename string
	mutex    sync.Mutex
	data     map[string]CacheItem
}

// NewFileCache creates a new cache instance
func NewFileCache(cacheName string) (*FileCache, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cacheDir := filepath.Join(homeDir, ".cache", "depshub")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	cacheFile := filepath.Join(cacheDir, cacheName+".json")
	cache := &FileCache{
		filename: cacheFile,
		data:     make(map[string]CacheItem),
	}

	// Load existing cache if it exists
	if err := cache.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return cache, nil
}

// Set adds or updates a cache entry with optional expiration duration
func (c *FileCache) Set(key string, value any, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Marshal the value to JSON bytes
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.data[key] = CacheItem{
		Value:      jsonBytes,
		CreateTime: time.Now(),
		ExpiresAt:  time.Now().Add(expiration),
	}

	return c.save()
}

// Get retrieves a value from the cache and unmarshals it into the provided destination
func (c *FileCache) Get(key string, dest any) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, exists := c.data[key]
	if !exists {
		return false, nil
	}

	// Check if item has expired
	if !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) {
		delete(c.data, key)
		err := c.save()
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// Unmarshal the JSON bytes into the destination
	if err := json.Unmarshal(item.Value, dest); err != nil {
		return true, err
	}

	return true, nil
}

// Delete removes an item from the cache
func (c *FileCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
	return c.save()
}

// Clear removes all items from the cache
func (c *FileCache) Clear() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]CacheItem)
	return c.save()
}

// save writes the cache to disk
func (c *FileCache) save() error {
	data, err := json.Marshal(c.data)
	if err != nil {
		return err
	}
	return os.WriteFile(c.filename, data, 0644)
}

// load reads the cache from disk
func (c *FileCache) load() error {
	data, err := os.ReadFile(c.filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.data)
}
