package netease

import (
	"strings"

	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

const (
	LinuxAPI          = "https://music.163.com/api/linux/forward"
	SearchAPI         = "https://music.163.com/weapi/search/get"
	GetSongAPI        = "https://music.163.com/weapi/v3/song/detail"
	GetSongURLAPI     = "https://music.163.com/weapi/song/enhance/player/url"
	GetArtistAPI      = "https://music.163.com/weapi/v1/artist/%d"
	GetAlbumAPI       = "https://music.163.com/weapi/v1/album/%d"
	GetPlaylistAPI    = "https://music.163.com/weapi/v3/playlist/detail"
	EmailLoginAPI     = "https://music.163.com/weapi/login"
	CellphoneLoginAPI = "https://music.163.com/weapi/login/cellphone"
	RefreshLoginAPI   = "https://music.163.com/weapi/login/token/refresh"
	LogoutAPI         = "https://music.163.com/weapi/logout"

	SongRequestLimit = 1000
	SongDefaultBR    = 128
)

var (
	std = New(provider.Client())
)

type (
	CommonResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg,omitempty"`
	}

	Song struct {
		Id      int      `json:"id"`
		Name    string   `json:"name"`
		Artists []Artist `json:"ar"`
		Album   Album    `json:"al"`
		Track   int      `json:"no"`
		Artist  string   `json:"-"`
		Lyric   string   `json:"-"`
		URL     string   `json:"-"`
	}

	SongSearchResponse struct {
		CommonResponse
		Result struct {
			Songs []*struct {
				Id      int    `json:"id"`
				Name    string `json:"name"`
				Artists []struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				} `json:"artists"`
				Album struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				} `json:"album"`
			} `json:"songs"`
			SongCount int `json:"songCount"`
		} `json:"result"`
	}

	SongURL struct {
		Id   int    `json:"id"`
		URL  string `json:"url"`
		BR   int    `json:"br"`
		Code int    `json:"code"`
	}

	SongResponse struct {
		CommonResponse
		Songs []*Song `json:"songs"`
	}

	SongURLResponse struct {
		CommonResponse
		Data []struct {
			Code int    `json:"code"`
			Id   int    `json:"id"`
			BR   int    `json:"br"`
			URL  string `json:"url"`
		} `json:"data"`
	}

	SongLyricResponse struct {
		CommonResponse
		Lrc struct {
			Lyric string `json:"lyric"`
		} `json:"lrc"`
		TLyric struct {
			Lyric string `json:"lyric"`
		} `json:"tlyric"`
	}

	Artist struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	}

	ArtistResponse struct {
		CommonResponse
		Artist struct {
			Artist
		} `json:"artist"`
		HotSongs []*Song `json:"hotSongs"`
	}

	Album struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	}

	AlbumResponse struct {
		CommonResponse
		Album Album   `json:"album"`
		Songs []*Song `json:"songs"`
	}

	PlaylistResponse struct {
		CommonResponse
		Playlist struct {
			Id       int     `json:"id"`
			Name     string  `json:"name"`
			PicURL   string  `json:"coverImgUrl"`
			Tracks   []*Song `json:"tracks"`
			TrackIds []struct {
				Id int `json:"id"`
			} `json:"trackIds"`
			Total int `json:"trackCount"`
		} `json:"playlist"`
	}

	LoginResponse struct {
		CommonResponse
		LoginType int `json:"loginType"`
		Account   struct {
			Id       int    `json:"id"`
			UserName string `json:"userName"`
		} `json:"account"`
	}

	API struct {
		Client *sreq.Client
	}
)

func (c *CommonResponse) String() string {
	return provider.ToJSON(c, false)
}

func (e *LoginResponse) String() string {
	return provider.ToJSON(e, false)
}

func (s *SongSearchResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongURLResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongLyricResponse) String() string {
	return provider.ToJSON(s, false)
}

func (a *ArtistResponse) String() string {
	return provider.ToJSON(a, false)
}

func (a *AlbumResponse) String() string {
	return provider.ToJSON(a, false)
}

func (p *PlaylistResponse) String() string {
	return provider.ToJSON(p, false)
}

func New(client *sreq.Client) *API {
	if client == nil {
		client = sreq.New(nil)
		cookie, _ := createCookie()
		client.SetDefaultRequestOpts(
			sreq.WithHeaders(sreq.Headers{
				"User-Agent": provider.UserAgent,
			}),
			sreq.WithCookies(cookie),
		)
	}
	return &API{
		Client: client,
	}
}

func Client() provider.API {
	return std
}

func (a *API) Platform() int {
	return provider.NetEase
}

func Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	return std.Request(method, url, opts...)
}

func (a *API) Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	cookie, _ := createCookie()
	defaultOpts := []sreq.RequestOption{
		sreq.WithHeaders(sreq.Headers{
			"Origin":  "https://music.163.com",
			"Referer": "https://music.163.com",
		}),
		sreq.WithCookies(cookie),
	}
	opts = append(opts, defaultOpts...)
	return a.Client.Request(method, url, opts...)
}

func (a *API) patchSongInfo(songs ...*Song) {
	for _, s := range songs {
		artists := make([]string, 0, len(s.Artists))
		for _, a := range s.Artists {
			artists = append(artists, strings.TrimSpace(a.Name))
		}
		s.Artist = strings.Join(artists, "/")
	}
}

func (a *API) patchSongURL(br int, songs ...*Song) {
	ids := make([]int, 0, len(songs))
	for _, s := range songs {
		ids = append(ids, s.Id)
	}

	resp, err := a.GetSongURLRaw(br, ids...)
	if err == nil && len(resp.Data) != 0 {
		m := make(map[int]string, len(resp.Data))
		for _, i := range resp.Data {
			m[i.Id] = i.URL
		}
		for _, s := range songs {
			s.URL = m[s.Id]
		}
	}
}

func (a *API) patchSongLyric(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			lyric, err := a.GetSongLyric(s.Id)
			if err == nil {
				s.Lyric = lyric
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) resolve(src []*Song) []*provider.Song {
	songs := make([]*provider.Song, 0, len(src))
	for _, s := range src {
		songs = append(songs, &provider.Song{
			Name:     strings.TrimSpace(s.Name),
			Artist:   s.Artist,
			Album:    strings.TrimSpace(s.Album.Name),
			PicURL:   s.Album.PicURL,
			Lyric:    s.Lyric,
			Playable: s.URL != "",
			URL:      s.URL,
		})
	}
	return songs
}
