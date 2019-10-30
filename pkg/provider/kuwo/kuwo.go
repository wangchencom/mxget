package kuwo

import (
	"net/http"
	"strings"

	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

const (
	SearchAPI        = "http://www.kuwo.cn/api/www/search/searchMusicBykeyWord"
	GetSongAPI       = "http://www.kuwo.cn/api/www/music/musicInfo"
	GetSongURLAPI    = "http://www.kuwo.cn/url?format=mp3&response=url&type=convert_url3"
	GetSongLyricAPI  = "http://www.kuwo.cn/newh5/singles/songinfoandlrc"
	GetArtistInfoAPI = "http://www.kuwo.cn/api/www/artist/artist"
	GetArtistSongAPI = "http://www.kuwo.cn/api/www/artist/artistMusic"
	GetAlbumAPI      = "http://www.kuwo.cn/api/www/album/albumInfo"
	GetPlaylistAPI   = "http://www.kuwo.cn/api/www/playlist/playListInfo"

	SongDefaultBR = 128
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
		RId             int    `json:"rid"`
		Name            string `json:"name"`
		ArtistId        int    `json:"artistid"`
		Artist          string `json:"artist"`
		AlbumId         int    `json:"albumid"`
		Album           string `json:"album"`
		AlbumPic        string `json:"albumpic"`
		Track           int    `json:"track"`
		IsListenFee     bool   `json:"isListenFee"`
		SongTimeMinutes string `json:"songTimeMinutes"`
		Lyric           string `json:"-"`
		URL             string `json:"-"`
	}

	SongSearchResponse struct {
		CommonResponse
		Data struct {
			Total string  `json:"total"`
			List  []*Song `json:"list"`
		} `json:"data"`
	}

	SongResponse struct {
		CommonResponse
		Data Song `json:"data"`
	}

	SongURLResponse struct {
		CommonResponse
		URL string `json:"url"`
	}

	SongLyricResponse struct {
		Status int    `json:"status"`
		Msg    string `json:"msg,omitempty"`
		Data   struct {
			LrcList []struct {
				Time      string `json:"time"`
				LineLyric string `json:"lineLyric"`
			} `json:"lrclist"`
		} `json:"data"`
	}

	ArtistInfo struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Pic300 string `json:"pic300"`
	}

	ArtistInfoResponse struct {
		CommonResponse
		Data ArtistInfo `json:"data"`
	}

	ArtistSongResponse struct {
		CommonResponse
		Data struct {
			List []*Song `json:"list"`
		} `json:"data"`
	}

	AlbumResponse struct {
		CommonResponse
		Data struct {
			AlbumId   int     `json:"albumId"`
			Album     string  `json:"album"`
			Pic       string  `json:"pic"`
			MusicList []*Song `json:"musicList"`
		} `json:"data"`
	}

	PlaylistResponse struct {
		CommonResponse
		Data struct {
			Id        int     `json:"id"`
			Name      string  `json:"name"`
			Img700    string  `json:"img700"`
			MusicList []*Song `json:"musicList"`
		} `json:"data"`
	}

	API struct {
		Client *sreq.Client
	}
)

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

func (a *ArtistInfoResponse) String() string {
	return provider.ToJSON(a, false)
}

func (a *ArtistSongResponse) String() string {
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

func (a *API) Platform() int {
	return provider.KuWo
}

func Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	return std.Request(method, url, opts...)
}

func (a *API) Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	// csrf 必须跟 kw_token 保持一致
	csrf := "0"
	cookie, err := a.Client.FilterCookie(url, "kw_token")
	if err == nil {
		csrf = cookie.Value
	}

	defaultOpts := []sreq.RequestOption{
		sreq.WithHeaders(sreq.Headers{
			"csrf":    csrf,
			"Origin":  "http://www.kuwo.cn",
			"Referer": "http://www.kuwo.cn",
		}),
	}
	if err != nil {
		defaultOpts = append(defaultOpts, sreq.WithCookies(
			&http.Cookie{
				Name:  "kw_token",
				Value: csrf,
			},
		), )
	}

	opts = append(opts, defaultOpts...)
	return a.Client.Request(method, url, opts...)
}

func (a *API) patchSongURL(br int, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			url, err := a.GetSongURL(s.RId, br)
			if err == nil {
				s.URL = url
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) patchSongLyric(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			lyric, err := a.GetSongLyric(s.RId)
			if err == nil {
				s.Lyric = lyric
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) resolve(src ...*Song) []*provider.Song {
	songs := make([]*provider.Song, 0, len(src))
	for _, s := range src {
		songs = append(songs, &provider.Song{
			Name:     strings.TrimSpace(s.Name),
			Artist:   strings.TrimSpace(strings.ReplaceAll(s.Artist, "&", "/")),
			Album:    strings.TrimSpace(s.Album),
			PicURL:   s.AlbumPic,
			Lyric:    s.Lyric,
			Playable: s.URL != "",
			URL:      s.URL,
		})
	}
	return songs
}
