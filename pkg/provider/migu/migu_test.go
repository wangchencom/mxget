package migu

import "testing"

func TestAPI_SearchSongs(t *testing.T) {
	result, err := SearchSongs("周杰伦")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestAPI_GetSong(t *testing.T) {
	song, err := GetSong("63273402938")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(song)
}

func TestAPI_GetSongURLRaw(t *testing.T) {
	resp, err := GetSongURLRaw("600908000002677565", "2")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := GetSongLyric("63273402938")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lyric)
}

func TestAPI_GetSongPic(t *testing.T) {
	pic, err := GetSongPic("1121439251")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pic)
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := GetArtist("112")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := GetAlbum("1121438701")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(album)
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := GetPlaylist("159248239")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(playlist)
}
