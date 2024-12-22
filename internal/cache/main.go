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

func InitializeCacheFile() (*cache.Cache, error) {
	var initErr error

	once.Do(func() {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			initErr = fmt.Errorf("failed to create cache directory: %w", err)
			return
		}

		if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
			file, err := os.Create(cacheFile)

			if err != nil {
				initErr = fmt.Errorf("failed to create cache file: %w", err)
				return
			}
			defer file.Close()

			if _, err := file.Write([]byte("{}")); err != nil {
				initErr = fmt.Errorf("failed to write to cache file: %w", err)
				return
			}
		}

		if err := LoadCacheFile(); err != nil {
			initErr = fmt.Errorf("failed to load cache file: %w", err)
			c = cache.New(5*time.Minute, 10*time.Minute)
			return
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return c, nil
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
