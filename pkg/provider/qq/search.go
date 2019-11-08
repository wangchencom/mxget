package qq

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func SearchSong(keyword string) (*provider.SearchSongsResult, error) {
	return std.SearchSong(keyword)
}

func (a *API) SearchSong(keyword string) (*provider.SearchSongsResult, error) {
	resp, err := a.SearchSongRaw(keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.Song.List)
	songs := make([]*provider.SearchSongsData, 0, n)
	for _, s := range resp.Data.Song.List {
		artists := make([]string, 0, len(s.Singer))
		for _, a := range s.Singer {
			artists = append(artists, strings.TrimSpace(a.Name))
		}
		songs = append(songs, &provider.SearchSongsData{
			Id:     s.Mid,
			Name:   strings.TrimSpace(s.Title),
			Artist: strings.Join(artists, "/"),
			Album:  strings.TrimSpace(s.Album.Name),
		})
	}
	return &provider.SearchSongsResult{
		Keyword: keyword,
		Count:   n,
		Songs:   songs,
	}, nil
}

func SearchSongRaw(keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	return std.SearchSongRaw(keyword, page, pageSize)
}

// 搜索歌曲
func (a *API) SearchSongRaw(keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	params := sreq.Params{
		"w": keyword,
		"p": strconv.Itoa(page),
		"n": strconv.Itoa(pageSize),
	}

	resp := new(SearchSongsResponse)
	err := a.Request(sreq.MethodGet, APISearch,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("search song: %d", resp.Code)
	}

	return resp, nil
}
