package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"

	"gopkg.in/yaml.v3"

	_ "github.com/joho/godotenv/autoload"
)

var DiscordToken = os.Getenv("BOT_TOKEN")
var SentryDSN = os.Getenv("SENTRY_DSN")
var Env = os.Getenv("GO_ENV")
var ConfigPath = os.Getenv("CONFIG_PATH")
var NinjaApiKey = os.Getenv("NINJA_API_KEY")

type ServerConfig struct {
	Id            string              `yaml:"id"`
	Name          string              `yaml:"name"`
	XpGiveEnabled bool                `yaml:"xpgive-enabled"`
	Roles         []RoleConfig        `yaml:"roles"`
	RatingRoles   []RatingRolesConfig `yaml:"ratings-roles"`
	PilotRoles    []RatingRolesConfig `yaml:"pilot-rating-roles"`
	Channels      []ChannelConfig     `yaml:"channels"`
	Emojis        []EmojiConfig       `yaml:"emojis"`
	BaseUrl       []BaseUrls          `yaml:"baseurl"`
	RoleRewards   []RoleReward        `yaml:"role_rewards"`
}

type RoleConfig struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type RatingRolesConfig struct {
	Name        string `yaml:"name"`
	RatingValue int    `yaml:"rating-value"`
	Id          string `yaml:"id"`
}

type ChannelConfig struct {
	Name        string               `yaml:"name"`
	Id          string               `yaml:"id"`
	Permissions []ChannelPermissions `yaml:"permissions"`
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

type ChannelPermissions struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

var Cfg ServerConfig

func GetXpGiveEnabled(id string) bool {
	cfg := configs[id]
	return cfg.XpGiveEnabled
}

func LoadAllServerConfigOrPanic(configPath string) map[string]ServerConfig {
	cfgs, err := LoadAllServerConfig(configPath)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf(err.Error())
	}
	return cfgs
}

func LoadAllServerConfig(configPath string) (map[string]ServerConfig, error) {
	cfgs := make(map[string]ServerConfig, 0)
	files, err := os.ReadDir(configPath)
	if err != nil {
		sentry.CaptureException(err)
		return nil, errors.New("failed to load server configs")
	}
	for _, f := range files {
		if !f.IsDir() {
			cfg, err := LoadServerConfig(fmt.Sprintf("%s/%s", configPath, f.Name()))
			if err != nil {
				sentry.CaptureException(err)
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
		sentry.CaptureException(err)
		return nil, err
	}
	var cfg ServerConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		sentry.CaptureException(err)
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
	return os.Getenv("INTERNAL_API_KEY")
}

// checks if the channel has the XP permission
func ValidXpChannel(id string, channel *discordgo.Channel) bool {
	channelName := channel.Name
	cfg := configs[id]

	if cfg.Channels == nil {
		return false
	}
	// Find the channel
	for i := 0; i < len(cfg.Channels); i++ {
		if strings.EqualFold(cfg.Channels[i].Name, channelName) {
			// Check the permission
			return GetBooleanPermissionValue(GetPermissionValue(cfg.Channels[i], "xp"))
		}
	}
	return false // return false either way - should we log this?
}

// Returns the permission value as a boolean based on what is passed in
// if the value is mispelled or not found it will return false
func GetBooleanPermissionValue(value string) bool {
	return strings.EqualFold(value, "true")
}

// Returns the permission value as a string
// All permission values are strings - this allows for greater flexibility -
// and handling of typos in the config (i.e. tuer or flsae admit it we have all done it)
// or additions of numerical values for a permission in the future
func GetPermissionValue(channel ChannelConfig, permissionName string) string {
	for i := 0; i < len(channel.Permissions); i++ {
		if strings.EqualFold(channel.Permissions[i].Name, permissionName) {
			return channel.Permissions[i].Value
		}
	}
	return "" // empty string if not found
}

func GetRatingsRoles(id string) []RatingRolesConfig {
	Cfg, _ := configs[id]
	return Cfg.RatingRoles
}

func GetPilotRatingsRoles(id string) []RatingRolesConfig {
	Cfg, _ := configs[id]
	return Cfg.PilotRoles
}
