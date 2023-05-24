package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileCache struct {
	cacheDir string
}

func NewFileCache(cacheDir string) (*FileCache, error) {
	// Use the OS temp folder if cacheDir is not provided
	if cacheDir == "" {
		cacheDir = os.TempDir()
	}

	// Create the cache directory if it doesn't exist
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &FileCache{cacheDir: cacheDir}, nil
}

func (fc *FileCache) Set(key string, value []byte, ttl int64) error {
	cacheFilePath := filepath.Join(fc.cacheDir, key)

	// Write the value to a file
	err := os.WriteFile(cacheFilePath, value, 0644)
	if err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// Set a timer to remove the file after the TTL expires
	if ttl > 0 {
		time.AfterFunc(time.Duration(ttl)*time.Second, func() {
			os.Remove(cacheFilePath)
		})
	}

	return nil
}

func (fc *FileCache) Get(key string) ([]byte, error) {
	cacheFilePath := filepath.Join(fc.cacheDir, key)

	// Read the value from the file
	value, err := os.ReadFile(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cache key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	return value, nil
}
