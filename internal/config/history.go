package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"workshop-downloader/internal/domain"
)

type HistoryManager struct {
	mu          sync.RWMutex
	historyPath string
	cfgManager  *ConfigManager
}

func NewHistoryManager(cfgManager *ConfigManager) *HistoryManager {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	appDir := filepath.Join(home, ".workshop-downloader")
	_ = os.MkdirAll(appDir, 0755)

	historyPath := filepath.Join(appDir, "history.json")
	return &HistoryManager{
		historyPath: historyPath,
		cfgManager:  cfgManager,
	}
}

func (hm *HistoryManager) LoadHistory() ([]domain.HistoryItem, error) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	data, err := os.ReadFile(hm.historyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.HistoryItem{}, nil
		}
		return nil, err
	}

	var items []domain.HistoryItem
	if err := json.Unmarshal(data, &items); err != nil {
		return []domain.HistoryItem{}, nil
	}

	// Update FolderExists status dynamically based on current OutputDir config
	cfg := hm.cfgManager.GetConfig()
	for i := range items {
		destDir := filepath.Join(cfg.OutputDir, items[i].AppID, items[i].WorkshopID)
		_, err := os.Stat(destDir)
		items[i].FolderExists = (err == nil)
		items[i].Path = destDir
	}

	return items, nil
}

func (hm *HistoryManager) SaveHistory(items []domain.HistoryItem) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(hm.historyPath, data, 0644)
}

func (hm *HistoryManager) AddHistoryItem(workshopID, appID, title, previewURL string) error {
	items, err := hm.LoadHistory()
	if err != nil {
		return err
	}

	// Avoid duplicates in history - update existing one or prepending
	existingIdx := -1
	for i, item := range items {
		if item.WorkshopID == workshopID && item.AppID == appID {
			existingIdx = i
			break
		}
	}

	newItem := domain.HistoryItem{
		ID:           fmt.Sprintf("hist_%d", time.Now().UnixNano()),
		WorkshopID:   workshopID,
		AppID:        appID,
		Title:        title,
		PreviewURL:   previewURL,
		DownloadedAt: time.Now(),
	}

	if existingIdx != -1 {
		// Update the existing item but move it to the front
		newItem.ID = items[existingIdx].ID
		items = append(items[:existingIdx], items[existingIdx+1:]...)
	}

	// Prepend to history so newest is first
	items = append([]domain.HistoryItem{newItem}, items...)

	return hm.SaveHistory(items)
}

func (hm *HistoryManager) DeleteHistoryItem(id string) error {
	items, err := hm.LoadHistory()
	if err != nil {
		return err
	}

	newItems := []domain.HistoryItem{}
	for _, item := range items {
		if item.ID != id {
			newItems = append(newItems, item)
		}
	}

	return hm.SaveHistory(newItems)
}
