package kuwo

import (
	"errors"
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
	resp, err := a.SearchSongsRaw(keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.List)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

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

func SearchSongsRaw(keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	return std.SearchSongsRaw(keyword, page, pageSize)
}

// 搜索歌曲
func (a *API) SearchSongsRaw(keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
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
			err = errors.New("search songs: no data")
		} else {
			err = fmt.Errorf("search songs: %s", resp.Msg)
		}
		return nil, err
	}

	return resp, nil
}
