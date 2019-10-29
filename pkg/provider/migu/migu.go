package migu

import (
	"strings"

	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

var (
	std = New(provider.Client())
)

const (
	SearchAPI        = "https://app.c.nf.migu.cn/MIGUM2.0/v1.0/content/search_all.do?isCopyright=1&isCorrect=1"
	GetSongIdAPI     = "http://music.migu.cn/v3/api/music/audioPlayer/songs?type=1"
	GetSongAPI       = "https://app.c.nf.migu.cn/MIGUM2.0/v2.0/content/querySongBySongId.do?contentId=0"
	GetSongURLAPI    = "https://app.c.nf.migu.cn/MIGUM2.0/v2.0/content/listen-url?copyrightId=0&netType=01&toneFlag=HQ"
	GetSongLyricAPI  = "http://music.migu.cn/v3/api/music/audioPlayer/getLyric"
	GetSongPicAPI    = "http://music.migu.cn/v3/api/music/audioPlayer/getSongPic"
	GetArtistInfoAPI = "https://app.c.nf.migu.cn/MIGUM2.0/v1.0/content/resourceinfo.do?needSimple=01&resourceType=2002"
	GetArtistSongAPI = "https://app.c.nf.migu.cn/MIGUM3.0/v1.0/template/singerSongs/release?templateVersion=2"
	GetAlbumAPI      = "https://app.c.nf.migu.cn/MIGUM2.0/v1.0/content/resourceinfo.do?needSimple=01&resourceType=2003"
	GetPlaylistAPI   = "https://app.c.nf.migu.cn/MIGUM2.0/v1.0/content/resourceinfo.do?needSimple=01&resourceType=2021"

	SongURL = "https://app.pd.nf.migu.cn/MIGUM2.0/v1.0/content/sub/listenSong.do?contentId=%s&copyrightId=0&netType=01&resourceType=%s&toneFlag=%s&channel=0"

	SongDefaultBR = 128
)

var codeRate = map[int]string{
	64:  "LQ",
	128: "PQ",
	320: "HQ",
	999: "SQ",
}

type (
	CommonResponse struct {
		Code string `json:"code"`
		Info string `json:"info,omitempty"`
	}

	SongSearchResult struct {
		ResourceType string `json:"resourceType"`
		ContentId    string `json:"contentId"`
		CopyrightId  string `json:"copyrightId"`
		Id           string `json:"id"`
		Name         string `json:"name"`
		Singers      []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"singers"`
		Albums []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"albums"`
	}

	SongSearchResponse struct {
		CommonResponse
		SongResultData struct {
			TotalCount string              `json:"totalCount"`
			Result     []*SongSearchResult `json:"result"`
		} `json:"songResultData"`
	}

	ImgItem struct {
		ImgSizeType string `json:"imgSizeType"`
		Img         string `json:"img"`
	}

	SongIdResponse struct {
		ReturnCode string `json:"returnCode"`
		Msg        string `json:"msg,omitempty"`
		Items      []struct {
			SongId string `json:"songId"`
		} `json:"items"`
	}

	Song struct {
		ResourceType string    `json:"resourceType"`
		ContentId    string    `json:"contentId"`
		CopyrightId  string    `json:"copyrightId"`
		SongId       string    `json:"songId"`
		SongName     string    `json:"songName"`
		SingerId     string    `json:"singerId"`
		Singer       string    `json:"singer"`
		AlbumId      string    `json:"albumId"`
		Album        string    `json:"album"`
		AlbumImgs    []ImgItem `json:"albumImgs"`
		LrcURL       string    `json:"lrcUrl"`
		Lyric        string    `json:"-"`
		URL          string    `json:"-"`
	}

	SongResponse struct {
		CommonResponse
		Resource []*Song `json:"resource"`
	}

	SongURLResponse struct {
		CommonResponse
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	SongLyricResponse struct {
		ReturnCode string `json:"returnCode"`
		Msg        string `json:"msg"`
		Lyric      string `json:"lyric"`
	}

	SongPicResponse struct {
		ReturnCode string `json:"returnCode"`
		Msg        string `json:"msg"`
		SmallPic   string `json:"smallPic"`
		MediumPic  string `json:"mediumPic"`
		LargePic   string `json:"largePic"`
	}

	ArtistInfo struct {
		ResourceType string    `json:"resourceType"`
		SingerId     string    `json:"singerId"`
		Singer       string    `json:"singer"`
		Imgs         []ImgItem `json:"imgs"`
	}

	ArtistInfoResponse struct {
		CommonResponse
		Resource []ArtistInfo `json:"resource"`
	}

	ArtistSongResponse struct {
		CommonResponse
		Data struct {
			ContentItemList []struct {
				ItemList []struct {
					Song Song `json:"song"`
				} `json:"itemList"`
			} `json:"contentItemList"`
		} `json:"data"`
	}

	Album struct {
		ResourceType string    `json:"resourceType"`
		AlbumId      string    `json:"albumId"`
		Title        string    `json:"title"`
		ImgItems     []ImgItem `json:"imgItems"`
		SongItems    []*Song   `json:"songItems"`
	}

	AlbumResponse struct {
		CommonResponse
		Resource []Album `json:"resource"`
	}

	Playlist struct {
		ResourceType string `json:"resourceType"`
		MusicListId  string `json:"musicListId"`
		Title        string `json:"title"`
		ImgItem      struct {
			Img string `json:"img"`
		} `json:"imgItem"`
		SongItems []*Song `json:"songItems"`
	}

	PlaylistResponse struct {
		CommonResponse
		Resource []Playlist `json:"resource"`
	}

	API struct {
		Client *sreq.Client
	}
)

func (s *SongSearchResponse) String() string {
	return provider.ToJSON(s, false)
}

func (s *SongIdResponse) String() string {
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

func (s *SongPicResponse) String() string {
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
	return provider.MiGu
}

func Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	return std.Request(method, url, opts...)
}

func (a *API) Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response {
	defaultOpts := []sreq.RequestOption{
		sreq.WithHeaders(sreq.Headers{
			"channel": "0",
			"Origin":  "http://music.migu.cn/v3",
			"Referer": "http://music.migu.cn/v3",
		}),
	}
	opts = append(opts, defaultOpts...)
	return a.Client.Request(method, url, opts...)
}

func (a *API) picURL(imgs []ImgItem) string {
	for _, i := range imgs {
		if i.ImgSizeType == "03" {
			return i.Img
		}
	}
	return ""
}

// 网页版API限流，并发请求经常503
// func (a *API) patchSongInfo(songs ...*Song) {
// 	c := concurrency.New(32)
// 	for _, s := range songs {
// 		c.Add(1)
// 		go func(s *Song) {
// 			picURL, err := a.GetSongPic(s.SongId)
// 			if err == nil {
// 				if !strings.HasPrefix(picURL, "http:") {
// 					picURL = "http:" + picURL
// 				}
// 				s.PicURL = picURL
// 			}
// 			c.Done()
// 		}(s)
// 	}
// 	c.Wait()
// }

func (a *API) patchSongURL(br int, songs ...*Song) {
	for _, s := range songs {
		s.URL = a.GetSongURL(s.ContentId, br)
	}
}

func (a *API) patchSongLyric(songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		c.Add(1)
		go func(s *Song) {
			lyric, err := a.Request(sreq.MethodGet, s.LrcURL).Text()
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
			Name:     strings.TrimSpace(s.SongName),
			Artist:   strings.TrimSpace(strings.ReplaceAll(s.Singer, "|", "/")),
			Album:    strings.TrimSpace(s.Album),
			PicURL:   a.picURL(s.AlbumImgs),
			Lyric:    s.Lyric,
			Playable: s.URL != "",
			URL:      s.URL,
		})
	}
	return songs
}
