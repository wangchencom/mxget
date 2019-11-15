package kugou

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/sreq"
)

func (a *API) SearchSongs(ctx context.Context, keyword string) (*api.SearchSongsResponse, error) {
	resp, err := a.SearchSongsRaw(ctx, keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.Info)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*api.SongResponse, n)
	for i, s := range resp.Data.Info {
		songs[i] = &api.SongResponse{
			Id:     s.Hash,
			Name:   strings.TrimSpace(s.SongName),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.SingerName, "、", "/")),
			Album:  strings.TrimSpace(s.AlbumName),
		}
	}
	return &api.SearchSongsResponse{
		Keyword: keyword,
		Count:   uint32(n),
		Songs:   songs,
	}, nil
}

// 搜索歌曲
func (a *API) SearchSongsRaw(ctx context.Context, keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	params := sreq.Params{
		"keyword":  keyword,
		"page":     strconv.Itoa(page),
		"pagesize": strconv.Itoa(pageSize),
	}

	resp := new(SearchSongsResponse)
	err := a.Request(sreq.MethodGet, APISearch,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("search songs: %s", resp.errorMessage())
	}

	return resp, nil
}
