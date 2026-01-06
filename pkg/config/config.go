package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	RoomID int    `json:"room_id"`
	Cookie string `json:"cookie"`
	Debug  bool   `json:"debug"`
}

var defaultConfig = &AppConfig{RoomID: 23990839, Cookie: "", Debug: false}

func GetConfigPath() string {
	userConfigDir, _ := os.UserConfigDir()
	appDir := filepath.Join(userConfigDir, "yuuna-danmu")

	_ = os.MkdirAll(appDir, 0o755)
	return filepath.Join(appDir, "config.json")
}

func Load() *AppConfig {
	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultConfig
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultConfig
	}
	return &cfg
}

func (cfg *AppConfig) Save() error {
	path := GetConfigPath()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
