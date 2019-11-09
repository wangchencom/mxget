package baidu

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

	n := len(resp.Result.SongInfo.SongList)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*provider.SearchSongsData, 0, n)
	for _, s := range resp.Result.SongInfo.SongList {
		songs = append(songs, &provider.SearchSongsData{
			Id:     s.SongId,
			Name:   strings.TrimSpace(s.Title),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.Author, ",", "/")),
			Album:  strings.TrimSpace(s.AlbumTitle),
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
		"query":     keyword,
		"page_no":   strconv.Itoa(page),
		"page_size": strconv.Itoa(pageSize),
	}

	resp := new(SearchSongsResponse)
	err := a.Request(sreq.MethodGet, APISearch,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("search songs: %s", resp.ErrorMessage)
	}

	return resp, nil
}
