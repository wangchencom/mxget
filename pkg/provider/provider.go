package provider

import (
	"bytes"
	"encoding/json"

	"github.com/winterssy/sreq"
)

const (
	NetEase = 1000 + iota
	QQ
	MiGu
	KuGou
	KuWo
	XiaMi
	BaiDu

	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
)

var (
	std *sreq.Client
)

type (
	Song struct {
		Name     string `json:"name"`
		Artist   string `json:"artist"`
		Album    string `json:"album"`
		PicURL   string `json:"pic_url,omitempty"`
		Lyric    string `json:"lyric,omitempty"`
		Playable bool   `json:"playable"`
		URL      string `json:"url,omitempty"`
	}

	Artist struct {
		Name   string  `json:"name"`
		PicURL string  `json:"pic_url,omitempty"`
		Count  int     `json:"count"`
		Songs  []*Song `json:"songs,omitempty"`
	}

	Album struct {
		Name   string  `json:"name"`
		PicURL string  `json:"pic_url,omitempty"`
		Count  int     `json:"count"`
		Songs  []*Song `json:"songs,omitempty"`
	}

	Playlist struct {
		Name   string  `json:"name"`
		PicURL string  `json:"pic_url,omitempty"`
		Count  int     `json:"count"`
		Songs  []*Song `json:"songs,omitempty"`
	}

	SearchSongsData struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
	}

	SearchSongsResult struct {
		Keyword string             `json:"keyword"`
		Count   int                `json:"count"`
		Songs   []*SearchSongsData `json:"songs,omitempty"`
	}

	Response struct {
		Code     int         `json:"code"`
		Msg      string      `json:"msg,omitempty"`
		Data     interface{} `json:"data,omitempty"`
		Platform int         `json:"platform,omitempty"`
	}

	API interface {
		// 平台标识
		Platform() int
		// 搜索歌曲
		SearchSongs(keyword string) (*SearchSongsResult, error)
		// 获取单曲
		GetSong(id string) (*Song, error)
		// 获取歌手
		GetArtist(id string) (*Artist, error)
		// 获取专辑
		GetAlbum(id string) (*Album, error)
		// 获取歌单
		GetPlaylist(id string) (*Playlist, error)
		// 网络请求
		Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response
	}
)

func init() {
	std = sreq.New(nil)
	std.SetDefaultRequestOpts(
		sreq.WithHeaders(sreq.Headers{
			"User-Agent": UserAgent,
		}),
	)
}

func Client() *sreq.Client {
	return std
}

func (s *SearchSongsResult) String() string {
	return ToJSON(s, false)
}

func (s *Song) String() string {
	return ToJSON(s, false)
}

func (a *Artist) String() string {
	return ToJSON(a, false)
}

func (a *Album) String() string {
	return ToJSON(a, false)
}

func (p *Playlist) String() string {
	return ToJSON(p, false)
}

func ToJSON(data interface{}, escapeHTML bool) string {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(escapeHTML)
	enc.SetIndent("", "\t")
	if enc.Encode(data) != nil {
		return "{}"
	}
	return buf.String()
}
