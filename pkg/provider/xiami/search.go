package xiami

import (
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func (a *API) SearchSongs(keyword string) (*provider.SearchSongsResult, error) {
	resp, err := a.SearchSongsRaw(keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.Data.Songs)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*provider.SearchSongsData, 0, n)
	for _, s := range resp.Data.Data.Songs {
		songs = append(songs, &provider.SearchSongsData{
			Id:     s.SongId,
			Name:   strings.TrimSpace(s.SongName),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.Singers, " / ", "/")),
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
