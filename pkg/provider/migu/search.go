package migu

import (
	"encoding/json"
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

	n := len(resp.SongResultData.Result)
	songs := make([]*provider.SearchSongsData, 0, n)
	for _, s := range resp.SongResultData.Result {
		artists := make([]string, 0, len(s.Singers))
		for _, a := range s.Singers {
			artists = append(artists, strings.TrimSpace(a.Name))
		}
		albums := make([]string, 0, len(s.Albums))
		for _, a := range s.Albums {
			albums = append(albums, strings.TrimSpace(a.Name))
		}
		songs = append(songs, &provider.SearchSongsData{
			Id:     s.CopyrightId,
			Name:   strings.TrimSpace(s.Name),
			Artist: strings.Join(artists, "/"),
			Album:  strings.Join(albums, "/"),
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
	switchOption := map[string]int{
		"song":     1,
		"album":    0,
		"singer":   0,
		"tagSong":  0,
		"mvSong":   0,
		"songlist": 0,
		"bestShow": 0,
	}
	enc, _ := json.Marshal(switchOption)
	params := sreq.Params{
		"searchSwitch": string(enc),
		"text":         keyword,
		"pageNo":       strconv.Itoa(page),
		"pageSize":     strconv.Itoa(pageSize),
	}

	resp := new(SearchSongsResponse)
	err := a.Request(sreq.MethodGet, APISearch,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "000000" {
		return nil, fmt.Errorf("search song: %s", resp.Info)
	}

	return resp, nil
}
