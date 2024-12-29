package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	_ "github.com/joho/godotenv/autoload"
)

var DiscordToken = os.Getenv("BOT_TOKEN")
var SentryDSN = os.Getenv("SENTRY_DSN")
var Env = os.Getenv("GO_ENV")
var ConfigPath = os.Getenv("CONFIG_PATH")

type ServerConfig struct {
	Id          string          `yaml:"id"`
	Name        string          `yaml:"name"`
	Roles       []RoleConfig    `yaml:"roles"`
	Channels    []ChannelConfig `yaml:"channels"`
	Emojis      []EmojiConfig   `yaml:"emojis"`
	BaseUrl     []BaseUrls      `yaml:"baseurl"`
	RoleRewards []RoleReward    `yaml:"role_rewards"`
	XpChannels  []XpChannel     `yaml:"xp_channels"`
}

type RoleConfig struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type ChannelConfig struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type EmojiConfig struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type BaseUrls struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
	Key  string `yaml:"key"`
}

type RoleReward struct {
	RoleName string `yaml:"role_name"`
	RoleID   string `yaml:"role_id"`
	Level    int    `yaml:"level"`
}

type XpChannel struct {
	Name string `yaml:"name"`
}

var Cfg ServerConfig

func LoadAllServerConfigOrPanic(configPath string) map[string]ServerConfig {
	cfgs, err := LoadAllServerConfig(configPath)
	if err != nil {
		log.Printf(err.Error())
	}
	return cfgs
}

func LoadAllServerConfig(configPath string) (map[string]ServerConfig, error) {
	cfgs := make(map[string]ServerConfig, 0)
	files, err := os.ReadDir(configPath)
	if err != nil {
		return nil, errors.New("failed to load server configs")
	}
	for _, f := range files {
		if !f.IsDir() {
			cfg, err := LoadServerConfig(fmt.Sprintf("%s/%s", configPath, f.Name()))
			if err != nil {
				log.Printf(err.Error())
				return nil, nil
			}
			cfgs[cfg.Id] = *cfg
		}
	}
	return cfgs, nil
}

func LoadServerConfig(configPath string) (*ServerConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg ServerConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	// TODO: Validate that roles aren't duplicated
	// TODO: Validate role criteria
	log.Printf("Loaded Config for %s (%s)\n", cfg.Name, cfg.Id)
	return &cfg, nil
}

var configs = LoadAllServerConfigOrPanic(ConfigPath)

func IntervalReloadConfigs() {
	for {
		time.Sleep(5 * time.Minute)
		log.Print("Reloading server configs")
		configs = LoadAllServerConfigOrPanic(ConfigPath)
	}
}

func GetServerConfig(id string) *ServerConfig {
	cfg, ok := configs[id]
	if !ok {
		return nil
	}
	return &cfg
}

func GetRoleId(id string, name string) string {
	cfg := configs[id]
	var RoleId string
	for i := 0; i < len(cfg.Roles); i++ {
		if cfg.Roles[i].Name == name {
			RoleId = cfg.Roles[i].Id
		}
	}
	return RoleId
}

func GetChannelId(id string, name string) string {
	cfg := configs[id]
	var ChannelId string
	for i := 0; i < len(cfg.Channels); i++ {
		if cfg.Channels[i].Name == name {
			ChannelId = cfg.Channels[i].Id
		}
	}
	return ChannelId
}

func GetEmojiId(id string, name string) string {
	cfg := configs[id]
	var EmojiId string
	for i := 0; i < len(cfg.Emojis); i++ {
		if cfg.Emojis[i].Name == name {
			EmojiId = cfg.Emojis[i].Id
		}
	}
	return EmojiId
}

func GetRoleRewards(id string) []RoleReward {
	cfg, exists := configs[id]
	if !exists {
		return nil
	}
	return cfg.RoleRewards
}

func GetApiBaseUrl(id string) string {
	cfg := configs[id]
	var BaseUrl string
	for i := 0; i < len(cfg.BaseUrl); i++ {
		if cfg.BaseUrl[i].Name == "Internal API" {
			BaseUrl = cfg.BaseUrl[i].Url
		}
	}
	return BaseUrl
}

func GetInternalApiKey(id string) string {
	cfg := configs[id]
	var ApiKey string
	for i := 0; i < len(cfg.BaseUrl); i++ {
		if cfg.BaseUrl[i].Name == "Internal API" {
			ApiKey = cfg.BaseUrl[i].Key
		}
	}
	return ApiKey
}

func ValidXpChannel(id string, channelName string) bool {
	cfg := configs[id]

	if cfg.XpChannels == nil {
		return false
	}

	for i := 0; i < len(cfg.XpChannels); i++ {
		if strings.EqualFold(cfg.XpChannels[i].Name, channelName) {
			return true
		}
	}
	return false
}
