package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	c          *cache.Cache
	once       sync.Once
	homeDir, _ = os.UserHomeDir()
	cacheFile  = filepath.Join(homeDir, ".eolctl", "cache.gob")
	cacheDir   = filepath.Dir(cacheFile)
)

func InitializeCacheFile() *cache.Cache {
	once.Do(func() {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			fmt.Printf("Failed to create cache directory: %v\n", err)
			return
		}

		if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
			file, err := os.Create(cacheFile)

			if err != nil {
				fmt.Println("Failed to create cache file")
				return
			}
			defer file.Close()

			if _, err := file.Write([]byte("{}")); err != nil {
				fmt.Printf("Failed to initialize cache file content: %v\n", err)
				return
			}
		}

		if err := LoadCacheFile(); err != nil {
			fmt.Printf("Failed to load cache file: %v\n", err)
			c = cache.New(5*time.Minute, 10*time.Minute)
			return
		}
	})

	return c
}

func LoadCacheFile() error {
	data, err := os.ReadFile(cacheFile)

	if err != nil {
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)

	var items map[string]cache.Item

	if err := dec.Decode(&items); err != nil {
		return fmt.Errorf("failed to decode cache file: %w", err)
	}

	c = cache.New(5*time.Minute, 10*time.Minute)

	for key, item := range items {
		c.Set(key, item, time.Duration(item.Expiration))
	}

	fmt.Println("Cache successfully loaded from cache file")
	return nil
}

func SaveCacheFile() error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(c.Items()); err != nil {
		return fmt.Errorf("failed to encode cache file: %w", err)
	}

	return os.WriteFile(cacheFile, buf.Bytes(), 0644)
}
