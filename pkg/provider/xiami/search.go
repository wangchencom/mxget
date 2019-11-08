package xiami

import (
	"fmt"
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

	n := len(resp.Data.Data.Songs)
	songs := make([]*provider.SearchSongsData, 0, n)
	for _, s := range resp.Data.Data.Songs {
		songs = append(songs, &provider.SearchSongsData{
			Id:     s.SongId,
			Name:   strings.TrimSpace(s.SongName),
			Artist: strings.TrimSpace(s.Singers),
			Album:  strings.TrimSpace(s.AlbumName),
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

func (a *API) SearchSongsRaw(keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	token, err := a.getToken(APISearch)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"key": keyword,
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	params := sreq.Params(signPayload(token, model))
	resp := new(SearchSongsResponse)
	err = a.Request(sreq.MethodGet, APISearch, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("search songs: %w", err)
	}

	return resp, nil
}
