package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/go-logr/logr"
)

const (
	cacheFileName = "hollowbeak/data.json"
)

type CacheItem struct {
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type Cache struct {
	logger    logr.Logger
	cacheFile string
	data      map[string]CacheItem
}

func NewCache(logger logr.Logger) (*Cache, error) {
	cacheFile, err := GetCachePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache path: %w", err)
	}

	cache := &Cache{
		logger:    logger,
		cacheFile: cacheFile,
		data:      make(map[string]CacheItem),
	}

	err = cache.load()
	if err != nil {
		return nil, fmt.Errorf("failed to load cache: %w", err)
	}

	return cache, nil
}

func (cache *Cache) Get(key string) (string, bool) {
	cache.logger.V(2).Info("Debug: Getting value from cache", "key", key)
	item, ok := cache.data[key]
	if !ok {
		return "", false
	}

	if time.Now().After(item.ExpiresAt) {
		cache.logger.V(2).Info("Debug: Cache item expired", "key", key)
		delete(cache.data, key)
		return "", false
	}

	return item.Value, true
}

func (cache *Cache) Set(key, value string) error {
	cache.logger.V(2).Info("Debug: Setting value in cache", "key", key)
	cache.data[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(6 * 30 * 24 * time.Hour),
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

func (cache *Cache) CleanupAndSave() error {
	cache.logger.V(1).Info("Debug: Starting cache cleanup")
	now := time.Now()
	for key, item := range cache.data {
		if now.After(item.ExpiresAt) {
			cache.logger.V(2).Info("Debug: Removing expired cache item", "key", key)
			delete(cache.data, key)
		}
	}
	cache.logger.V(1).Info("Debug: Cache cleanup completed", "remainingEntries", len(cache.data))

	return cache.save()
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

func GetCachePath() (string, error) {
	configFilePath, err := xdg.ConfigFile(cacheFileName)
	if err != nil {
		return "", fmt.Errorf("failed to get XDG config file path: %w", err)
	}
	return filepath.Abs(configFilePath)
}
