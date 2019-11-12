package qq

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/mxget/pkg/utils"
	"github.com/winterssy/sreq"
)

const (
	APISearch        = "https://c.y.qq.com/soso/fcgi-bin/client_search_cp?format=json&platform=yqq&new_json=1"
	APIGetSong       = "https://c.y.qq.com/v8/fcg-bin/fcg_play_single_song.fcg?format=json&platform=yqq"
	APIGetSongURLV1  = "http://c.y.qq.com/base/fcgi-bin/fcg_music_express_mobile3.fcg?format=json&platform=yqq&needNewCode=0&cid=205361747&uin=0&guid=0"
	APIGetSongsURLV2 = "https://u.y.qq.com/cgi-bin/musicu.fcg?format=json&platform=yqq"
	APIGetSongLyric  = "https://c.y.qq.com/lyric/fcgi-bin/fcg_query_lyric_new.fcg?format=json&platform=yqq&nobase64=1"
	APIGetArtist     = "https://c.y.qq.com/v8/fcg-bin/fcg_v8_singer_track_cp.fcg?format=json&platform=yqq&newsong=1&order=listen"
	APIGetAlbum      = "https://c.y.qq.com/v8/fcg-bin/fcg_v8_album_detail_cp.fcg?format=json&platform=yqq&newsong=1"
	APIGetPlaylist   = "https://c.y.qq.com/v8/fcg-bin/fcg_v8_playlist_cp.fcg?format=json&platform=yqq&newsong=1"

	SongURL      = "http://mobileoc.music.tc.qq.com/%s?guid=0&uin=0&vkey=%s"
	ArtistPicURL = "https://y.gtimg.cn/music/photo_new/T001R800x800M000%s.jpg"
	AlbumPicURL  = "https://y.gtimg.cn/music/photo_new/T002R800x800M000%s.jpg"

	SongURLRequestLimit = 300
)

var (
	std = New(provider.Client())
)

type (
	CommonResponse struct {
		Code int `json:"code"`
	}

	Song struct {
		Mid    string   `json:"mid"`
		Title  string   `json:"title"`
		Singer []Singer `json:"singer"`
		Album  Album    `json:"album"`
		Track  int      `json:"index_album"`
		Action struct {
			Switch int `json:"switch"`
		} `json:"action"`
		File struct {
			MediaMid string `json:"media_mid"`
		} `json:"file"`
		Lyric string `json:"-"`
		URL   string `json:"-"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Data struct {
			Song struct {
				TotalNum int     `json:"totalnum"`
				List     []*Song `json:"list"`
			} `json:"song"`
		} `json:"data"`
	}

	SongResponse struct {
		CommonResponse
		Data []*Song `json:"data"`
	}

	SongURLResponseV1 struct {
		Code    int    `json:"code"`
		Cid     int    `json:"cid"`
		ErrInfo string `json:"errinfo,omitempty"`
		Data    struct {
			Expiration int `json:"expiration"`
			Items      []struct {
				SubCode  int    `json:"subcode"`
				SongMid  string `json:"songmid"`
				FileName string `json:"filename"`
				Vkey     string `json:"vkey"`
			} `json:"items"`
		} `json:"data"`
	}

	SongURLResponseV2 struct {
		CommonResponse
		Req0 struct {
			Data struct {
				MidURLInfo []struct {
					FileName string `json:"filename"`
					PURL     string `json:"purl"`
					SongMid  string `json:"songmid"`
					Vkey     string `json:"vkey"`
				} `json:"midurlinfo"`
				Sip        []string `json:"sip"`
				TestFile2g string   `json:"testfile2g"`
			} `json:"data"`
		} `json:"req0"`
	}

	SongLyricResponse struct {
		CommonResponse
		Lyric string `json:"lyric"`
		Trans string `json:"trans"`
	}

	Singer struct {
		Mid  string `json:"mid"`
		Name string `json:"name"`
	}

	ArtistResponse struct {
		CommonResponse
		Data struct {
			SingerMid  string `json:"singer_mid"`
			SingerName string `json:"singer_name"`
			List       []struct {
				MusicData *Song `json:"musicData"`
			} `json:"list"`
		} `json:"data"`
	}

	Album struct {
		Mid  string `json:"mid"`
		Name string `json:"name"`
	}

	AlbumResponse struct {
		CommonResponse
		Data struct {
			GetAlbumInfo struct {
				FAlbumMid  string `json:"Falbum_mid"`
				FAlbumName string `json:"Falbum_name"`
			} `json:"getAlbumInfo"`
			GetSongInfo []*Song `json:"getSongInfo"`
		} `json:"data"`
	}

	PlaylistResponse struct {
		CommonResponse
		Data struct {
			CDList []struct {
				DissTid  string  `json:"disstid"`
				DissName string  `json:"dissname"`
				Logo     string  `json:"logo"`
				PicURL   string  `json:"dir_pic_url2"`
				SongList []*Song `json:"songlist"`
			} `json:"cdlist"`
		} `json:"data"`
	}

	API struct {
		Client *sreq.Client
	}
)

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

func (s *SongResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongURLResponseV2) String() string {
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

func (a *API) Platform() int {
	return provider.QQ
}

func (a *API) Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	defaultOpts := []sreq.RequestOption{
		sreq.WithHeaders(sreq.Headers{
			"Origin":  "https://c.y.qq.com",
			"Referer": "https://c.y.qq.com",
		}),
	}
	opts = append(opts, defaultOpts...)
	return a.Client.Request(method, url, opts...)
}

func (a *API) patchSongURLV1(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			url, err := a.GetSongURLV1(s.Mid, s.File.MediaMid)
			if err == nil {
				s.URL = url
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) patchSongURLV2(songs ...*Song) {
	n := len(songs)
	songMids := make([]string, 0, n)
	for _, s := range songs {
		songMids = append(songMids, s.Mid)
	}

	type result struct {
		resp *SongURLResponseV2
		err  error
	}

	urlMap := make(map[string]string, n)
	queue := make(chan *result)
	wg := new(sync.WaitGroup)
	// url长度限制，每次请求的歌曲数不能太多，分批获取
	for i := 0; i < n; i += SongURLRequestLimit {
		ids := songMids[i:utils.Min(i+SongURLRequestLimit, n)]
		wg.Add(1)
		go func() {
			resp, err := a.GetSongsURLV2Raw(ids...)
			queue <- &result{
				resp: resp,
				err:  err,
			}
		}()
	}
	go func() {
		for r := range queue {
			if r.err == nil {
				n := len(r.resp.Req0.Data.Sip)
				if n > 0 {
					// 随机获取一个sip
					sip := r.resp.Req0.Data.Sip[rand.Intn(n)]
					for _, i := range r.resp.Req0.Data.MidURLInfo {
						if i.PURL != "" {
							urlMap[i.SongMid] = sip + i.PURL
						}
					}
				}
			}
			wg.Done()
		}
	}()
	wg.Wait()

	for _, s := range songs {
		s.URL = urlMap[s.Mid]
	}
}

func (a *API) patchSongLyric(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			lyric, err := a.GetSongLyric(s.Mid)
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
		artists := make([]string, 0, len(s.Singer))
		for _, a := range s.Singer {
			artists = append(artists, strings.TrimSpace(a.Name))
		}
		songs = append(songs, &provider.Song{
			Name:     strings.TrimSpace(s.Title),
			Artist:   strings.Join(artists, "/"),
			Album:    strings.TrimSpace(s.Album.Name),
			PicURL:   fmt.Sprintf(AlbumPicURL, s.Album.Mid),
			Lyric:    s.Lyric,
			Playable: s.URL != "",
			URL:      s.URL,
		})
	}
	return songs
}
