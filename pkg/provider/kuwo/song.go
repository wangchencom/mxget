package kuwo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetSong(mid string) (*provider.Song, error) {
	return std.GetSong(mid)
}

func (a *API) GetSong(mid string) (*provider.Song, error) {
	resp, err := a.GetSongRaw(mid)
	if err != nil {
		return nil, err
	}

	a.patchSongURL(SongDefaultBR, &resp.Data)
	a.patchSongLyric(&resp.Data)
	songs := resolve(&resp.Data)
	return songs[0], nil
}

func GetSongRaw(mid string) (*SongResponse, error) {
	return std.GetSongRaw(mid)
}

// 获取歌曲详情
func (a *API) GetSongRaw(mid string) (*SongResponse, error) {
	params := sreq.Params{
		"mid": mid,
	}

	resp := new(SongResponse)
	err := a.Request(sreq.MethodGet, APIGetSong,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get song: %s", resp.Msg)
	}

	return resp, nil
}

func GetSongURL(mid int, br int) (string, error) {
	return std.GetSongURL(mid, br)
}

func (a *API) GetSongURL(mid int, br int) (string, error) {
	resp, err := a.GetSongURLRaw(mid, br)
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}

func GetSongURLRaw(mid int, br int) (*SongURLResponse, error) {
	return std.GetSongURLRaw(mid, br)
}

// 获取歌曲播放地址，br: 比特率，128/192/320 可选
func (a *API) GetSongURLRaw(mid int, br int) (*SongURLResponse, error) {
	var _br int
	switch br {
	case 128, 192, 320:
		_br = br
	default:
		_br = 320
	}
	params := sreq.Params{
		"rid": strconv.Itoa(mid),
		"br":  fmt.Sprintf("%dkmp3", _br),
	}

	resp := new(SongURLResponse)
	err := a.Request(sreq.MethodGet, APIGetSongURL,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get song url: %s", resp.Msg)
	}

	return resp, nil
}

func GetSongLyric(mid int) (string, error) {
	return std.GetSongLyric(mid)
}

func (a *API) GetSongLyric(mid int) (string, error) {
	resp, err := a.GetSongLyricRaw(mid)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, i := range resp.Data.LrcList {
		t, err := strconv.ParseFloat(i.Time, 64)
		if err != nil {
			return "", err
		}
		m := int(t / 60)
		s := int(t - float64(m*60))
		ms := int((t - float64(m*60) - float64(s)) * 100)
		sb.WriteString(fmt.Sprintf("[%02d:%02d:%02d]%s\n", m, s, ms, i.LineLyric))
	}

	return sb.String(), nil
}

func GetSongLyricRaw(mid int) (*SongLyricResponse, error) {
	return std.GetSongLyricRaw(mid)
}

// 获取歌词
func (a *API) GetSongLyricRaw(mid int) (*SongLyricResponse, error) {
	params := sreq.Params{
		"musicId": strconv.Itoa(mid),
	}

	resp := new(SongLyricResponse)
	err := a.Request(sreq.MethodGet, APIGetSongLyric,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("get song lyric: %s", resp.Msg)
	}

	return resp, nil
}
