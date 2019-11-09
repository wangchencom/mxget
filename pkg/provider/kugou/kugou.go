package kugou

import (
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

const (
	APISearch           = "http://mobilecdn.kugou.com/api/v3/search/song"
	APIGetSong          = "http://m.kugou.com/app/i/getSongInfo.php?cmd=playInfo"
	APIGetSongURL       = "http://trackercdn.kugou.com/i/v2/?pid=2&behavior=play&cmd=25"
	APIGetSongLyric     = "http://m.kugou.com/app/i/krc.php?cmd=100&timelength=1"
	APIGetArtistInfo    = "http://mobilecdn.kugou.com/api/v3/singer/info"
	APIGetArtistSongs   = "http://mobilecdn.kugou.com/api/v3/singer/song"
	APIGetAlbumInfo     = "http://mobilecdn.kugou.com/api/v3/album/info"
	APIGetAlbumSongs    = "http://mobilecdn.kugou.com/api/v3/album/song"
	APIGetPlaylistInfo  = "http://mobilecdn.kugou.com/api/v3/special/info"
	APIGetPlaylistSongs = "http://mobilecdn.kugou.com/api/v3/special/song"
)

var (
	std = New(provider.Client())
)

type (
	CommonResponse struct {
		Status  int    `json:"status"`
		Error   string `json:"error,omitempty"`
		ErrCode int    `json:"errcode"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Data struct {
			Total int `json:"total"`
			Info  []*struct {
				Hash       string `json:"hash"`
				HQHash     string `json:"320hash"`
				SQHash     string `json:"sqhash"`
				SongName   string `json:"songname"`
				SingerName string `json:"singername"`
				AlbumId    string `json:"album_id"`
				AlbumName  string `json:"album_name"`
			} `json:"info"`
		} `json:"data"`
	}

	Song struct {
		Hash         string `json:"hash"`
		SongName     string `json:"songName"`
		SingerId     int    `json:"singerId"`
		SingerName   string `json:"singerName"`
		ChoricSinger string `json:"choricSinger"`
		FileName     string `json:"fileName"`
		ExtName      string `json:"extName"`
		AlbumId      int    `json:"albumid"`
		AlbumImg     string `json:"album_img"`
		Extra        struct {
			PQHash string `json:"128hash"`
			HQHash string `json:"320hash"`
			SQHash string `json:"sqhash"`
		} `json:"extra"`
		URL       string `json:"url"`
		AlbumName string `json:"-"`
		Lyric     string `json:"-"`
	}

	SongResponse struct {
		CommonResponse
		Song
	}

	SongURLResponse struct {
		Status  int      `json:"status"`
		BitRate int      `json:"bitRate"`
		ExtName string   `json:"extName"`
		URL     []string `json:"url"`
	}

	ArtistInfo struct {
		SingerId   int    `json:"singerid"`
		SingerName string `json:"singername"`
		ImgURL     string `json:"imgurl"`
	}

	ArtistInfoResponse struct {
		CommonResponse
		Data ArtistInfo `json:"data"`
	}

	ArtistSongsResponse struct {
		CommonResponse
		Data struct {
			Info []*Song `json:"info"`
		} `json:"data"`
	}

	AlbumInfo struct {
		AlbumId   int    `json:"albumid"`
		AlbumName string `json:"albumname"`
		ImgURL    string `json:"imgurl"`
	}

	AlbumInfoResponse struct {
		CommonResponse
		Data AlbumInfo `json:"data"`
	}

	AlbumSongsResponse struct {
		CommonResponse
		Data struct {
			Info []*Song `json:"info"`
		} `json:"data"`
	}

	PlaylistInfo struct {
		SpecialId   int    `json:"specialid"`
		SpecialName string `json:"specialname"`
		ImgURL      string `json:"imgurl"`
	}

	PlaylistInfoResponse struct {
		CommonResponse
		Data PlaylistInfo `json:"data"`
	}

	PlaylistSongsResponse struct {
		CommonResponse
		Data struct {
			Info []*Song `json:"info"`
		} `json:"data"`
	}

	API struct {
		Client *sreq.Client
	}
)

func (s *SearchSongsResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongURLResponse) String() string {
	return provider.ToJSON(s, false)
}

func (a *ArtistInfoResponse) String() string {
	return provider.ToJSON(a, false)
}

func (a *ArtistSongsResponse) String() string {
	return provider.ToJSON(a, false)
}

func (a *AlbumInfoResponse) String() string {
	return provider.ToJSON(a, false)
}

func (a *AlbumSongsResponse) String() string {
	return provider.ToJSON(a, false)
}

func (p *PlaylistInfoResponse) String() string {
	return provider.ToJSON(p, false)
}

func (p *PlaylistSongsResponse) String() string {
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
	return provider.KuGou
}

func Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	return std.Request(method, url, opts...)
}

func (a *API) Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	defaultOpts := []sreq.RequestOption{
		sreq.WithHeaders(sreq.Headers{
			"Origin":  "https://www.kugou.com",
			"Referer": "https://www.kugou.com",
		}),
	}
	opts = append(opts, defaultOpts...)
	return a.Client.Request(method, url, opts...)
}

func (a *API) patchSongInfo(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			resp, err := a.GetSongRaw(s.Hash)
			if err == nil {
				s.SongName = resp.SongName
				s.SingerId = resp.SingerId
				s.SingerName = resp.SingerName
				s.ChoricSinger = resp.ChoricSinger
				s.AlbumId = resp.AlbumId
				s.AlbumImg = resp.AlbumImg
				s.Extra = resp.Extra
				s.URL = resp.URL
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) patchSongURL(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if s.URL != "" {
			continue
		}
		c.Add(1)
		go func(s *Song) {
			url, err := a.GetSongURL(s.Hash)
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
			lyric, err := a.GetSongLyric(s.Hash)
			if err == nil {
				s.Lyric = lyric
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) patchAlbumInfo(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			if s.AlbumId != 0 {
				resp, err := a.GetAlbumInfoRaw(strconv.Itoa(s.AlbumId))
				if err == nil {
					s.AlbumName = resp.Data.AlbumName
				}
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func resolve(src ...*Song) []*provider.Song {
	songs := make([]*provider.Song, 0, len(src))
	for _, s := range src {
		songs = append(songs, &provider.Song{
			Name:     strings.TrimSpace(s.SongName),
			Artist:   strings.TrimSpace(strings.ReplaceAll(s.ChoricSinger, "、", "/")),
			Album:    strings.TrimSpace(s.AlbumName),
			PicURL:   strings.ReplaceAll(s.AlbumImg, "{size}", "480"),
			Lyric:    s.Lyric,
			Playable: s.URL != "",
			URL:      s.URL,
		})
	}
	return songs
}
