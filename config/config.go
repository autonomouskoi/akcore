package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kirsle/configdir"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AutonomousKoi *AutonomousKoi `yaml:"autonomouskoi"`
	Autoban       *Autoban       `yaml:"autoban"`
	Dummies       *Dummies       `yaml:"dummies"`
	IRC           *IRC           `yaml:"irc"`
	Log           *Log           `yaml:"log"`
	Magic         *Magic         `yaml:"magic"`
	Modules       *Modules       `yaml:"modules"`
	OSCIn         *OSCIn         `yaml:"oscin"`
	Plugins       *Plugins       `yaml:"plugins"`
	Spam          *Spam          `yaml:"spam"`
	Twitch        *Twitch        `yaml:"twitch"`
	Web           *Web           `yaml:"web"`
}

func (c *Config) relativize(base string) {
	for _, subCfg := range []interface {
		relativize(string)
	}{
		c.Autoban,
		c.Dummies,
		c.IRC,
		c.Log,
		c.Magic,
		c.OSCIn,
		c.Plugins,
		c.Spam,
		c.Twitch,
		c.Web,
	} {
		subCfg.relativize(base)
	}
}

type AutonomousKoi struct {
	StoragePath string
	CachePath   string
}

type Autoban struct {
	BotsPath string `yaml:"bots_path"`
	Enabled  bool   `yaml:"enabled"`
}

func (a *Autoban) relativize(base string) {
	relativizePath(base, &a.BotsPath)
}

type Dummies struct {
	Count int    `yaml:"count"`
	Path  string `yaml:"path"`
}

func (d *Dummies) relativize(base string) {
	relativizePath(base, &d.Path)
}

type IRC struct {
	Nick         string   `yaml:"nick"`
	ServerAddr   string   `yaml:"server_addr"`
	ServerPort   string   `yaml:"server_port"`
	Capabilities []string `yaml:"capabilities"`
	SoftCmdsPath string   `yaml:"soft_cmds_path"`
}

func (i *IRC) relativize(base string) {
	relativizePath(base, &i.SoftCmdsPath)
}

type Log struct {
	Dir string `yaml:"dir"`
}

func (l *Log) relativize(base string) {
	relativizePath(base, &l.Dir)
}

type Magic struct {
	AvatarPath         string             `yaml:"avatar_path"`
	KappaPath          string             `yaml:"kappa_path"`
	MagicOSCHost       string             `yaml:"magic_osc_host"`
	MagicOSCPort       int                `yaml:"magic_osc_port"`
	OverlayMessagePath string             `yaml:"overlay_msg_path"`
	OverlayTitlePath   string             `yaml:"overlay_title_path"`
	Overlays           map[string]float64 `yaml:"overlays"`
	Underlays          map[string]float64 `yaml:"underlays"`

	SceneIDs     map[string]float64 `yaml:"scene_ids"`
	LockedScenes []string           `yaml:"locked_scenes"`
}

func (c *Magic) relativize(base string) {
	relativizePath(base, &c.AvatarPath)
	relativizePath(base, &c.KappaPath)
	relativizePath(base, &c.OverlayMessagePath)
	relativizePath(base, &c.OverlayTitlePath)
}

type Modules struct {
	Disabled StringSet `yaml:"disabled"`
}

type Plugins struct {
	Path string `yaml:"path"`
}

func (p *Plugins) relativize(base string) {
	relativizePath(base, &p.Path)
}

type OSCIn struct {
	Addr string `yaml:"addr"`
}

func (*OSCIn) relativize(base string) {}

type Spam struct {
	Interval time.Duration     `yaml:"interval"`
	Cooldown time.Duration     `yaml:"cooldown"`
	Spams    map[string]string `yaml:"spams"`
}

func (*Spam) relativize(base string) {}

type Twitch struct {
	AvatarCachePath      string `yaml:"avatar_cache_path"`
	BroadcasterID        string `yaml:"broadcaster_id"`
	BroadcasterTokenPath string `yaml:"broadcaster_token_path"`
	BotTokenPath         string `yaml:"bot_token_path"`
	ClientID             string `yaml:"client_id"`
	ClientSecret         string `yaml:"client_secret"`
}

func (t *Twitch) relativize(base string) {
	relativizePath(base, &t.AvatarCachePath)
	relativizePath(base, &t.BroadcasterTokenPath)
	relativizePath(base, &t.BotTokenPath)
}

type Web struct {
	CachePath  string `yaml:"cache_path"`
	Listen     string `yaml:"listen"`
	StaticPath string `yaml:"static_path"`
}

func (w *Web) relativize(base string) {
	relativizePath(base, &w.CachePath)
	relativizePath(base, &w.StaticPath)
}

func New(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}
	dir := filepath.Dir(path)
	cfg.relativize(dir)
	cfg.AutonomousKoi = &AutonomousKoi{
		StoragePath: filepath.Join(configdir.LocalConfig("autonomouskoi")),
		CachePath:   filepath.Join(configdir.LocalCache("autonomouskoi")),
	}
	return &cfg, nil
}

func relativizePath(base string, path *string) {
	if filepath.IsAbs(*path) {
		return
	}
	*path = filepath.Join(base, *path)
}
