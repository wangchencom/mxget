package qq

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

	n := len(resp.Data.Song.List)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*provider.SearchSongsData, n)
	for i, s := range resp.Data.Song.List {
		artists := make([]string, len(s.Singer))
		for j, a := range s.Singer {
			artists[j] = strings.TrimSpace(a.Name)
		}
		songs[i] = &provider.SearchSongsData{
			Id:     s.Mid,
			Name:   strings.TrimSpace(s.Title),
			Artist: strings.Join(artists, "/"),
			Album:  strings.TrimSpace(s.Album.Name),
		}
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
		return nil, fmt.Errorf("search songs: %d", resp.Code)
	}

	return resp, nil
}
