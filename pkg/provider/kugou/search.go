package kugou

import (
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

	n := len(resp.Data.Info)
	songs := make([]*provider.SearchSongData, 0, n)
	for _, s := range resp.Data.Info {
		songs = append(songs, &provider.SearchSongData{
			Id:     s.Hash,
			Name:   strings.TrimSpace(s.SongName),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.SingerName, "、", "/")),
			Album:  strings.TrimSpace(s.AlbumName),
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
		"keyword":  keyword,
		"page":     strconv.Itoa(page),
		"pagesize": strconv.Itoa(pageSize),
	}

	resp := new(SongSearchResponse)
	err := a.Request(sreq.MethodGet, SearchAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("search song: %s", resp.Error)
	}

	return resp, nil
}
