package settings

import (
	"encoding/json"
	"errors"
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
	configFileName       = "mxget.json"
	hiddenConfigFileName = "." + configFileName
	downloadDir          = "./downloads"
)

var (
	CfgPath = makeConfigPath()
	Cfg     *Config
	Limit   int
	Tag     bool
	Lyric   bool
	Force   bool

	errorConfigFileNotFound = errors.New("config file not found")
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
	}
)

func (c *Config) Save() error {
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(CfgPath, b, 0644)
}

func Load() {
	conf, err := load()
	if err != nil {
		if err == errorConfigFileNotFound {
			conf = &Config{
				DownloadDir:   downloadDir,
				MusicPlatform: provider.NetEase,
			}
			err = conf.Save()
			if err != nil {
				easylog.Errorf("Failed to make config file: %v", err)
			}
		} else {
			easylog.Fatalf("Failed to load config file: %v", err)
		}
	}

	if !VerifyPlatform(conf.MusicPlatform) {
		easylog.Errorf("Unexpected music platform: %d", conf.MusicPlatform)
		easylog.Infof("Reset to default")
		conf.MusicPlatform = provider.NetEase
		_ = conf.Save()
	}

	if err := os.MkdirAll(conf.DownloadDir, 0755); err != nil {
		easylog.Errorf("Failed to make download dir: %v", err)
		easylog.Infof("Reset to default")
		conf.DownloadDir = downloadDir
		_ = conf.Save()
	}
	Cfg = conf
}

func load() (*Config, error) {
	b, err := ioutil.ReadFile(CfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errorConfigFileNotFound
		}
		return nil, err
	}

	conf := new(Config)
	err = json.Unmarshal(b, conf)
	return conf, err
}

func makeConfigPath() string {
	var cfgDir string
	var cfgPath string

	homeDir, err := os.UserHomeDir()
	if err == nil {
		cfgDir = filepath.Join(homeDir, ".config", "mxget")
		cfgPath = filepath.Join(cfgDir, configFileName)
	}

	_, err = os.Stat(cfgPath)
	if err == nil {
		return cfgPath
	}

	if cfgPath != "" {
		err := os.MkdirAll(cfgDir, 0755)
		if err == nil {
			return cfgPath
		}
	}

	return hiddenConfigFileName
}

func Platform(flag string) int {
	return platform[flag]
}

func Client(platform int) provider.API {
	return client[platform]
}

func Site(platform int) string {
	return site[platform]
}

func VerifyPlatform(platform int) bool {
	switch platform {
	case provider.NetEase, provider.QQ, provider.MiGu, provider.KuGou, provider.KuWo:
		return true
	default:
		return false
	}
}
