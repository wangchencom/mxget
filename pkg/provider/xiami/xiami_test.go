package xiami

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
	song, err := GetSong("xMPr7Lbbb28")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(song)
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := GetSongLyric("xMPr7Lbbb28")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lyric)
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := GetArtist("3110")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := GetAlbum("nmTM4c70144")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(album)
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := GetPlaylist("8007523")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(playlist)
}
