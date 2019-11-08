package netease

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func SearchSongs(keyword string) (*provider.SearchSongsResult, error) {
	return std.SearchSongs(keyword)
}

func (a *API) SearchSongs(keyword string) (*provider.SearchSongsResult, error) {
	resp, err := a.SearchSongsRaw(keyword, 0, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Result.Songs)
	songs := make([]*provider.SearchSongsData, 0, n)
	for _, s := range resp.Result.Songs {
		artists := make([]string, 0, len(s.Artists))
		for _, a := range s.Artists {
			artists = append(artists, strings.TrimSpace(a.Name))
		}
		songs = append(songs, &provider.SearchSongsData{
			Id:     strconv.Itoa(s.Id),
			Name:   strings.TrimSpace(s.Name),
			Artist: strings.Join(artists, "/"),
			Album:  s.Album.Name,
		})
	}
	return &provider.SearchSongsResult{
		Keyword: keyword,
		Count:   n,
		Songs:   songs,
	}, nil
}

func SearchSongsRaw(keyword string, offset int, limit int) (*SearchSongsResponse, error) {
	return std.SearchSongsRaw(keyword, offset, limit)
}

// 搜索歌曲
func (a *API) SearchSongsRaw(keyword string, offset int, limit int) (*SearchSongsResponse, error) {
	// type: 1: 单曲, 10: 专辑, 100: 歌手, 1000: 歌单, 1002: 用户,
	// 1004: MV, 1006: 歌词, 1009: 电台, 1014: 视频
	data := map[string]interface{}{
		"s":      keyword,
		"type":   1,
		"offset": offset,
		"limit":  limit,
	}

	resp := new(SearchSongsResponse)
	err := a.Request(sreq.MethodPost, APISearch,
		sreq.WithForm(weapi(data)),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("search songs: %s", resp.Msg)
	}

	return resp, nil
}
