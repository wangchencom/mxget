package qq

import "testing"

func TestAPI_SearchSongs(t *testing.T) {
	result, err := SearchSongs("Alan Walker")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestAPI_GetSong(t *testing.T) {
	song, err := GetSong("002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(song)
}

func TestAPI_GetSongURL(t *testing.T) {
	url, err := GetSongURL("002Zkt5S2z8JZx", "002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := GetSongLyric("002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lyric)
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := GetArtist("000Sp0Bz4JXH0o")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := GetAlbum("002fRO0N4FftzY")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(album)
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := GetPlaylist("5474239760")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(playlist)
}
