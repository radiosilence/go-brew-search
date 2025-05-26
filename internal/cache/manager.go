package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Manager struct {
	dir string
	ttl time.Duration
}

type cacheEntry struct {
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

func New(dir string, ttl time.Duration) *Manager {
	return &Manager{
		dir: dir,
		ttl: ttl,
	}
}

func (m *Manager) Get(key string, dest interface{}) error {
	path := m.cachePath(key)
	
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cache miss: %s", key)
		}
		return err
	}
	
	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return err
	}
	
	// Check if cache is expired
	if time.Since(entry.Timestamp) > m.ttl {
		os.Remove(path)
		return fmt.Errorf("cache expired: %s", key)
	}
	
	// Marshal to JSON and unmarshal to destination
	jsonData, err := json.Marshal(entry.Data)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(jsonData, dest)
}

func (m *Manager) Set(key string, data interface{}) error {
	entry := cacheEntry{
		Data:      data,
		Timestamp: time.Now(),
	}
	
	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	
	path := m.cachePath(key)
	return os.WriteFile(path, jsonData, 0644)
}

func (m *Manager) Clear() error {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			os.Remove(filepath.Join(m.dir, entry.Name()))
		}
	}
	
	return nil
}

func (m *Manager) cachePath(key string) string {
	return filepath.Join(m.dir, key+".json")
}