package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)
import _ "github.com/joho/godotenv/autoload"

var DiscordToken = os.Getenv("BOT_TOKEN")
var SentryDSN = os.Getenv("SENTRY_DSN")
var Env = os.Getenv("GO_ENV")
var ConfigPath = os.Getenv("CONFIG_PATH")

type ServerConfig struct {
	Id       string          `yaml:"id"`
	Name     string          `yaml:"name"`
	Roles    []RoleConfig    `yaml:"roles"`
	Channels []ChannelConfig `yaml:"channels"`
	Emojis   []EmojiConfig   `yaml:"emojis"`
	BaseUrl  []BaseUrls      `yaml:"baseurl"`
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
	Cfg, _ := configs[id]
	var RoleId string
	for i := 0; i < len(Cfg.Roles); i++ {
		if Cfg.Roles[i].Name == name {
			RoleId = Cfg.Roles[i].Id
		}
	}
	return RoleId
}

func GetChannelId(id string, name string) string {
	Cfg, _ := configs[id]
	var ChannelId string
	for i := 0; i < len(Cfg.Channels); i++ {
		if Cfg.Channels[i].Name == name {
			ChannelId = Cfg.Channels[i].Id
		}
	}
	return ChannelId
}

func GetEmojiId(id string, name string) string {
	Cfg, _ := configs[id]
	var EmojiId string
	for i := 0; i < len(Cfg.Emojis); i++ {
		if Cfg.Emojis[i].Name == name {
			EmojiId = Cfg.Emojis[i].Id
		}
	}
	return EmojiId
}
