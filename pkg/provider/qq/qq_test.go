package qq_test

import (
	"os"
	"testing"

	"github.com/winterssy/mxget/pkg/provider/qq"
)

var client *qq.API

func setup() {
	client = qq.New(nil)
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestAPI_SearchSongs(t *testing.T) {
	result, err := client.SearchSongs("Alan Walker")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestAPI_GetSong(t *testing.T) {
	song, err := client.GetSong("002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(song)
}

func TestAPI_GetSongURLV1(t *testing.T) {
	url, err := client.GetSongURLV1("002Zkt5S2z8JZx", "002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}

func TestAPI_GetSongURLV2(t *testing.T) {
	url, err := client.GetSongURLV2("002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := client.GetSongLyric("002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lyric)
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := client.GetArtist("000Sp0Bz4JXH0o")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := client.GetAlbum("002fRO0N4FftzY")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(album)
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := client.GetPlaylist("5474239760")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(playlist)
}
