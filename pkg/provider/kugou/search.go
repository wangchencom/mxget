package kugou

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func (a *API) SearchSongs(keyword string) (*provider.SearchSongsResult, error) {
	resp, err := a.SearchSongsRaw(keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.Info)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*provider.SearchSongsData, 0, n)
	for _, s := range resp.Data.Info {
		songs = append(songs, &provider.SearchSongsData{
			Id:     s.Hash,
			Name:   strings.TrimSpace(s.SongName),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.SingerName, "、", "/")),
			Album:  strings.TrimSpace(s.AlbumName),
		})
	}
	return &provider.SearchSongsResult{
		Keyword: keyword,
		Count:   n,
		Songs:   songs,
	}, nil
}

// 搜索歌曲
func (a *API) SearchSongsRaw(keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	params := sreq.Params{
		"keyword":  keyword,
		"page":     strconv.Itoa(page),
		"pagesize": strconv.Itoa(pageSize),
	}

	resp := new(SearchSongsResponse)
	err := a.Request(sreq.MethodGet, APISearch,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("search songs: %s", resp.Error)
	}

	return resp, nil
}
