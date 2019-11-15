package baidu

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/sreq"
)

func (a *API) GetArtist(ctx context.Context, tingUid string) (*api.ArtistResponse, error) {
	resp, err := a.GetArtistRaw(ctx, tingUid, 0, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.SongList)
	if n == 0 {
		return nil, errors.New("get artist: no data")
	}

	a.patchSongsURL(ctx, resp.SongList...)
	a.patchSongsLyric(ctx, resp.SongList...)
	songs := resolve(resp.SongList...)
	return &api.ArtistResponse{
		Id:     resp.ArtistInfo.TingUid,
		Name:   strings.TrimSpace(resp.ArtistInfo.Name),
		PicUrl: strings.SplitN(resp.ArtistInfo.AvatarBig, "@", 2)[0],
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取歌手
func (a *API) GetArtistRaw(ctx context.Context, tingUid string, offset int, limits int) (*ArtistResponse, error) {
	params := sreq.Params{
		"tinguid": tingUid,
		"offset":  strconv.Itoa(offset),
		"limits":  strconv.Itoa(limits),
	}

	resp := new(ArtistResponse)
	err := a.Request(sreq.MethodGet, APIGetArtist,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get artist: %v", resp.errorMessage())
	}

	return resp, nil
}
