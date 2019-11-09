package baidu

import (
	"testing"
)

func TestAPI_SearchSongs(t *testing.T) {
	result, err := SearchSongs("五月天")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestAPI_GetSong(t *testing.T) {
	song, err := GetSong("1686649")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(song)
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := GetSongLyric("1686649")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lyric)
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := GetArtist("1557")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := GetAlbum("946499")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(album)
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := GetPlaylist("566347665")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(playlist)
}
