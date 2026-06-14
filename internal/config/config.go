package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"workshop-downloader/internal/domain"
)

type ConfigManager struct {
	mu         sync.RWMutex
	configPath string
	config     domain.AppConfig
}

func NewConfigManager() *ConfigManager {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	appDir := filepath.Join(home, ".workshop-downloader")
	_ = os.MkdirAll(appDir, 0755)

	configPath := filepath.Join(appDir, "config.json")
	cm := &ConfigManager{
		configPath: configPath,
		config: domain.AppConfig{
			SteamCmdPath:   filepath.Join(appDir, "steamcmd", "steamcmd.exe"),
			OutputDir:      filepath.Join(home, "Downloads", "SteamWorkshop"),
			AutoUpdate:     true,
			Username:       "anonymous",
			MaxConcurrency: 2,
		},
	}
	_ = cm.Load()
	return cm
}

func (cm *ConfigManager) GetConfig() domain.AppConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

func (cm *ConfigManager) UpdateConfig(cfg domain.AppConfig) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config = cfg
	return cm.Save()
}

func (cm *ConfigManager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		// File might not exist yet, use default paths
		return nil
	}

	err = json.Unmarshal(data, &cm.config)
	if err == nil {
		if cm.config.MaxConcurrency <= 0 {
			cm.config.MaxConcurrency = 2
		}
	}
	return err
}

func (cm *ConfigManager) Save() error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cm.configPath, data, 0644)
}
