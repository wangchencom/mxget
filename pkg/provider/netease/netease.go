package netease

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

const (
	APILinux          = "https://music.163.com/api/linux/forward"
	APISearch         = "https://music.163.com/weapi/search/get"
	APIGetSongs       = "https://music.163.com/weapi/v3/song/detail"
	APIGetSongsURL    = "https://music.163.com/weapi/song/enhance/player/url"
	APIGetArtist      = "https://music.163.com/weapi/v1/artist/%d"
	APIGetAlbum       = "https://music.163.com/weapi/v1/album/%d"
	APIGetPlaylist    = "https://music.163.com/weapi/v3/playlist/detail"
	APIEmailLogin     = "https://music.163.com/weapi/login"
	APICellphoneLogin = "https://music.163.com/weapi/login/cellphone"
	APIRefreshLogin   = "https://music.163.com/weapi/login/token/refresh"
	APILogout         = "https://music.163.com/weapi/logout"

	SongRequestLimit = 1000
	SongDefaultBR    = 128
)

var (
	std = New(provider.Client())

	cookie *http.Cookie
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
		Lyric   string   `json:"-"`
		URL     string   `json:"-"`
	}

	SearchSongsResponse struct {
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

	SongsResponse struct {
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

func init() {
	cookie, _ = createCookie()
}

func New(client *sreq.Client) *API {
	if client == nil {
		client = sreq.New(nil)
		client.SetDefaultRequestOpts(
			sreq.WithHeaders(sreq.Headers{
				"User-Agent": provider.UserAgent,
			}),
		)
	}
	return &API{
		Client: client,
	}
}

func Client() provider.API {
	return std
}

func (s *SearchSongsResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongsResponse) String() string {
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

func (e *LoginResponse) String() string {
	return provider.ToJSON(e, false)
}

func (a *API) PlatformId() provider.PlatformId {
	return provider.NetEase
}

func (a *API) Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	defaultOpts := []sreq.RequestOption{
		sreq.WithHeaders(sreq.Headers{
			"Origin":  "https://music.163.com",
			"Referer": "https://music.163.com",
		}),
	}

	// 如果已经登录，不需要额外设置cookie，cookie jar会自动管理
	_, err := a.Client.FilterCookie(url, "MUSIC_U")
	if err != nil {
		defaultOpts = append(defaultOpts, sreq.WithCookies(cookie))
	}

	opts = append(opts, defaultOpts...)
	return a.Client.Request(method, url, opts...)
}

func (a *API) patchSongURL(br int, songs ...*Song) {
	ids := make([]int, 0, len(songs))
	for _, s := range songs {
		ids = append(ids, s.Id)
	}

	resp, err := a.GetSongsURLRaw(br, ids...)
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

func resolve(src ...*Song) []*provider.Song {
	songs := make([]*provider.Song, 0, len(src))
	for _, s := range src {
		artists := make([]string, 0, len(s.Artists))
		for _, a := range s.Artists {
			artists = append(artists, strings.TrimSpace(a.Name))
		}
		songs = append(songs, &provider.Song{
			Id:       strconv.Itoa(s.Id),
			Name:     strings.TrimSpace(s.Name),
			Artist:   strings.Join(artists, "/"),
			Album:    strings.TrimSpace(s.Album.Name),
			PicURL:   s.Album.PicURL,
			Lyric:    s.Lyric,
			Playable: s.URL != "",
			URL:      s.URL,
		})
	}
	return songs
}
