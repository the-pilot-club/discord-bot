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
var FCPToken = os.Getenv("FCP_TOKEN")
var FCPEnv = os.Getenv("FCP_ENV")
var CoreAPIToken = os.Getenv("INTERNAL_API_KEY")
var RedisURL = os.Getenv("REDIS_URL")

type ServerConfig struct {
	Id                   string              `yaml:"id"`
	Name                 string              `yaml:"name"`
	XpGiveEnabled        bool                `yaml:"xpgive-enabled"`
	EventReminders       bool                `yaml:"event-reminders"`
	Roles                []RoleConfig        `yaml:"roles"`
	RatingRoles          []RatingRolesConfig `yaml:"ratings-roles"`
	PilotRoles           []RatingRolesConfig `yaml:"pilot-rating-roles"`
	Channels             []ChannelConfig     `yaml:"channels"`
	XpDisabledChannels   []ChannelReference  `yaml:"xp-disabled-channels"`
	XpDisabledCategories []ChannelReference  `yaml:"xp-disabled-categories"`
	Emojis               []EmojiConfig       `yaml:"emojis"`
	BaseUrl              []BaseUrls          `yaml:"baseurl"`
	RoleRewards          []RoleReward        `yaml:"role_rewards"`
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
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

// ChannelReference points at a channel or category. It can be written in the
// config either as the Discord ID or as a mapping with name
// and id, e.g.
//
//	xp-disabled-channels:
//	  - 1234567890
//	  - name: Staff Chat
//	    id: 9876543210
type ChannelReference struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

func (r *ChannelReference) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		r.Id = value.Value
		return nil
	}
	// avoid recursing into this method when decoding the mapping form
	type rawReference ChannelReference
	var raw rawReference
	if err := value.Decode(&raw); err != nil {
		return err
	}
	*r = ChannelReference(raw)
	return nil
}

type EmojiConfig struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type BaseUrls struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
	URL  string `yaml:"url"`
	Key  string `yaml:"key"`
}

func (b BaseUrls) Value() string {
	if strings.TrimSpace(b.Link) != "" {
		return b.Link
	}
	return b.URL
}

type RoleReward struct {
	RoleName string `yaml:"role_name"`
	RoleID   string `yaml:"role_id"`
	Level    int    `yaml:"level"`
}

var Cfg ServerConfig

func GetXpGiveEnabled(id string) bool {
	cfg := configs[id]
	return cfg.XpGiveEnabled
}
func EventRemindersEnabled(id string) bool {
	cfg := configs[id]
	return cfg.EventReminders
}

func LoadAllServerConfigOrPanic(configPath string) map[string]ServerConfig {
	cfgs, err := LoadAllServerConfig(configPath)
	if err != nil {
		sentry.CaptureException(err)
		log.Print(err.Error())
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
				log.Print(err.Error())
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

func GetBaseUrl(id string, name string) string {
	cfg := configs[id]
	var BaseUrl string
	for i := 0; i < len(cfg.BaseUrl); i++ {
		if cfg.BaseUrl[i].Name == name {
			BaseUrl = cfg.BaseUrl[i].Value()
		}
	}
	return BaseUrl
}

func GetInternalApiKey(id string) string {
	return os.Getenv("INTERNAL_API_KEY")
}

// ValidXpChannel reports whether messages in the given channel should earn XP.
// XP is earned everywhere by default. A guild opts specific channels or whole
// categories out via xp-disabled-channels / xp-disabled-categories in its
// config. Threads inherit the state of trheir parent channel
func ValidXpChannel(guildID string, channel *discordgo.Channel, ancestors ...*discordgo.Channel) bool {
	if channel == nil {
		return false
	}

	cfg := configs[guildID]

	for _, ch := range append([]*discordgo.Channel{channel}, ancestors...) {
		if ch == nil {
			continue
		}
		if channelMatchesRef(cfg.XpDisabledChannels, ch) || channelMatchesRef(cfg.XpDisabledCategories, ch) {
			return false
		}
		// A parent we weren't handed an object for can still be matched by ID
		if idInRefs(cfg.XpDisabledChannels, ch.ParentID) || idInRefs(cfg.XpDisabledCategories, ch.ParentID) {
			return false
		}
	}

	return true
}

// reports whether channel is named in refs, by ID or name
func channelMatchesRef(refs []ChannelReference, channel *discordgo.Channel) bool {
	if channel == nil {
		return false
	}
	normalizedName := normalizeChannelName(channel.Name)
	for _, ref := range refs {
		if ref.Id != "" && ref.Id == channel.ID {
			return true
		}
		if ref.Name != "" && normalizeChannelName(ref.Name) == normalizedName {
			return true
		}
	}
	return false
}

// idInRefs reports whether id matches the ID of any reference in refs
func idInRefs(refs []ChannelReference, id string) bool {
	if id == "" {
		return false
	}
	for _, ref := range refs {
		if ref.Id == id {
			return true
		}
	}
	return false
}

func normalizeChannelName(name string) string {
	replacer := strings.NewReplacer(" ", "-", "_", "-")
	return strings.ToLower(replacer.Replace(name))
}

func GetRatingsRoles(id string) []RatingRolesConfig {
	Cfg, _ := configs[id]
	return Cfg.RatingRoles
}

func GetPilotRatingsRoles(id string) []RatingRolesConfig {
	Cfg, _ := configs[id]
	return Cfg.PilotRoles
}
