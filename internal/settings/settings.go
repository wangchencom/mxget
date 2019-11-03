package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/winterssy/easylog"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/mxget/pkg/provider/kugou"
	"github.com/winterssy/mxget/pkg/provider/kuwo"
	"github.com/winterssy/mxget/pkg/provider/migu"
	"github.com/winterssy/mxget/pkg/provider/netease"
	"github.com/winterssy/mxget/pkg/provider/qq"
)

const (
	downloadDir = "./downloads"
)

var (
	Cfg = &Config{
		DownloadDir:   downloadDir,
		MusicPlatform: provider.NetEase,
	}
	Limit int
	Tag   bool
	Lyric bool
	Force bool
)

var (
	platform = map[string]int{
		"netease": provider.NetEase,
		"nc":      provider.NetEase,
		"qq":      provider.QQ,
		"migu":    provider.MiGu,
		"mg":      provider.MiGu,
		"kugou":   provider.KuGou,
		"kg":      provider.KuGou,
		"kuwo":    provider.KuGou,
		"kw":      provider.KuWo,
	}

	client = map[int]provider.API{
		provider.NetEase: netease.Client(),
		provider.QQ:      qq.Client(),
		provider.MiGu:    migu.Client(),
		provider.KuGou:   kugou.Client(),
		provider.KuWo:    kuwo.Client(),
	}

	site = map[int]string{
		provider.NetEase: "music.163.com",
		provider.QQ:      "y.qq.com",
		provider.MiGu:    "music.migu.cn",
		provider.KuGou:   "kugou.com",
		provider.KuWo:    "kuwo.cn",
	}
)

type (
	Config struct {
		DownloadDir   string `json:"download_dir"`
		MusicPlatform int    `json:"music_platform"`

		others   map[string]interface{} `json:"-"`
		filePath string                 `json:"-"`
	}
)

func GetPlatformId(platformFlag string) int {
	return platform[platformFlag]
}

func GetClient(platformId int) provider.API {
	return client[platformId]
}

func GetSite(platformId int) string {
	return site[platformId]
}

func Init() {
	err := Cfg.setup()
	if err != nil {
		Cfg.Reset()
		easylog.Fatalf("Can't initialize settings, reset to defaults: %v", err)
	}
}

func (c *Config) setup() error {
	c.getCfgFile()
	err := c.loadCfgFile()
	if err != nil {
		return err
	}

	err = c.check()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) getCfgFile() {
	var cfgDir string
	xdgDir := os.Getenv("XDG_CONFIG_HOME")
	if xdgDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			cfgDir = "."
		} else {
			cfgDir = filepath.Join(home, ".config", "mxget")
		}
	} else {
		cfgDir = filepath.Join(xdgDir, "mxget")
	}

	if os.MkdirAll(cfgDir, 0755) != nil {
		c.filePath = ".mxget.json"
	}

	c.filePath = filepath.Join(cfgDir, "mxget.json")
}

func (c *Config) loadCfgFile() error {
	_, err := os.Stat(c.filePath)
	if err == nil {
		b, err := ioutil.ReadFile(c.filePath)
		if err != nil {
			return err
		}
		return json.Unmarshal(b, c)
	}

	return c.Save()
}

func (c *Config) check() error {
	if GetSite(c.MusicPlatform) == "" {
		c.MusicPlatform = provider.NetEase
		return fmt.Errorf("unexpected music platform: %d", c.MusicPlatform)
	}

	err := os.MkdirAll(c.DownloadDir, 0755)
	if err != nil {
		c.DownloadDir = downloadDir
		return fmt.Errorf("cant't make download dir: %w", err)
	}

	return nil
}

func (c *Config) Save() error {
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.filePath, b, 0644)
}

func (c *Config) Reset() {
	_ = c.Save()
}
