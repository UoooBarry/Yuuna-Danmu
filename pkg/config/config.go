package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	RoomID       int              `json:"room_id"`
	Cookie       string           `json:"cookie"`
	RefreshToken string           `json:"refresh_token"`
	Debug        bool             `json:"debug"`
	Servers      []ServerSettings `json:"servers"`
	Transparent  bool             `json:"transparent"`
}

type ServerSettings struct {
	Name    string     `json:"name"`
	Type    ServerType `json:"type"`
	Port    int        `json:"port"`
	Enabled bool       `json:"enabled"`
}

var defaultConfig = &AppConfig{
	RoomID: 23990839,
	Cookie: "",
	Debug:  false,
	Servers: []ServerSettings{
		{Name: "gRPC", Type: "grpc", Port: 50051, Enabled: false},
	},
	RefreshToken: "",
	Transparent:  true,
}

type ServerType string

var GRPC ServerType = "grpc"

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
	if cfg.Servers == nil {
		cfg.Servers = defaultConfig.Servers
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
