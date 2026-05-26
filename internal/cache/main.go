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

			enc := gob.NewEncoder(file)
			if err := enc.Encode(map[string]cache.Item{}); err != nil {
				initErr = fmt.Errorf("failed to initialize cache file: %w", err)
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
		var ttl time.Duration
		if item.Expiration > 0 {
			ttl = time.Until(time.Unix(0, item.Expiration))
			if ttl <= 0 {
				continue // item already expired
			}
		} else {
			ttl = cache.NoExpiration
		}
		c.Set(key, item.Object, ttl)
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
