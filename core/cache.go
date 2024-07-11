package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/go-logr/logr"
)

type Cache struct {
	logger    logr.Logger
	cacheFile string
	data      map[string]string
}

func NewCache(logger logr.Logger) (*Cache, error) {
	cacheFile, err := getCachePath("hollowbeak/data.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get cache path: %w", err)
	}

	cache := &Cache{
		logger:    logger,
		cacheFile: cacheFile,
		data:      make(map[string]string),
	}

	err = cache.load()
	if err != nil {
		return nil, fmt.Errorf("failed to load cache: %w", err)
	}

	return cache, nil
}

func (cache *Cache) Get(key string) (string, bool) {
	cache.logger.V(2).Info("Debug: Getting value from cache", "key", key)
	value, ok := cache.data[key]
	return value, ok
}

func (cache *Cache) Set(key, value string) error {
	cache.logger.V(2).Info("Debug: Setting value in cache", "key", key)
	cache.data[key] = value
	err := cache.save()
	if err != nil {
		cache.logger.Error(err, "Failed to save cache after setting value", "key", key)
		return fmt.Errorf("failed to save cache after setting value: %w", err)
	}
	return nil
}

func (cache *Cache) load() error {
	cache.logger.V(1).Info("Debug: Loading cache", "path", cache.cacheFile)
	data, err := os.ReadFile(cache.cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			cache.logger.Info("Cache file not found, starting with empty cache", "path", cache.cacheFile)
			return nil
		}
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	err = json.Unmarshal(data, &cache.data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	cache.logger.V(1).Info("Debug: Cache loaded successfully", "entries", len(cache.data))
	return nil
}

func (cache *Cache) save() error {
	cache.logger.V(1).Info("Debug: Saving cache", "path", cache.cacheFile)
	data, err := json.MarshalIndent(cache.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	err = os.WriteFile(cache.cacheFile, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	cache.logger.V(1).Info("Debug: Cache saved successfully")
	return nil
}

func getCachePath(configRelPath string) (string, error) {
	configFilePath, err := xdg.ConfigFile(configRelPath)
	if err != nil {
		return "", fmt.Errorf("failed to get XDG config file path: %w", err)
	}

	dirPerm := os.FileMode(0o700)
	dir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return configFilePath, nil
}
