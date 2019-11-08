package kuwo

import (
	"errors"
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

	n := len(resp.Data.List)
	songs := make([]*provider.SearchSongsData, 0, n)
	for _, s := range resp.Data.List {
		songs = append(songs, &provider.SearchSongsData{
			Id:     strconv.Itoa(s.RId),
			Name:   strings.TrimSpace(s.Name),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.Artist, "&", "/")),
			Album:  strings.TrimSpace(s.Album),
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
		"key": keyword,
		"pn":  strconv.Itoa(page),
		"rn":  strconv.Itoa(pageSize),
	}

	resp := new(SearchSongsResponse)
	err := a.Request(sreq.MethodGet, APISearch,
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
