package song

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SongDetail struct {
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func FetchSongDetails(group, song string) (*SongDetail, error) {
	url := fmt.Sprintf("http://music-api-service/info?group=%s&song=%s", group, song)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var details SongDetail
	err = json.NewDecoder(resp.Body).Decode(&details)
	return &details, err
}
