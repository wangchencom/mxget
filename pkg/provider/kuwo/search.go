package kuwo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func SearchSong(keyword string) (*provider.SearchResult, error) {
	return std.SearchSong(keyword)
}

func (a *API) SearchSong(keyword string) (*provider.SearchResult, error) {
	resp, err := a.SearchSongRaw(keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.List)
	songs := make([]*provider.SearchSongData, 0, n)
	for _, s := range resp.Data.List {
		songs = append(songs, &provider.SearchSongData{
			Id:     strconv.Itoa(s.RId),
			Name:   strings.TrimSpace(s.Name),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.Artist, "&", "/")),
			Album:  strings.TrimSpace(s.Album),
		})
	}
	return &provider.SearchResult{
		Keyword: keyword,
		Count:   n,
		Songs:   songs,
	}, nil
}

func SearchSongRaw(keyword string, page int, pageSize int) (*SongSearchResponse, error) {
	return std.SearchSongRaw(keyword, page, pageSize)
}

// 搜索歌曲
func (a *API) SearchSongRaw(keyword string, page int, pageSize int) (*SongSearchResponse, error) {
	params := sreq.Params{
		"key": keyword,
		"pn":  strconv.Itoa(page),
		"rn":  strconv.Itoa(pageSize),
	}

	resp := new(SongSearchResponse)
	err := a.Request(sreq.MethodGet, SearchAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		if resp.Code == -1 {
			err = errors.New("search song: no data")
		} else {
			err = fmt.Errorf("search song: %s", resp.Msg)
		}
		return nil, err
	}

	return resp, nil
}
