package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"song-library/internal/config"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func FetchSongInfo(group, song string) (*SongDetail, error) {
	baseURL := config.AppConfig.API.ExternalURL + "/info"
	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Отправляем запрос
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch song info: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Декодируем JSON-ответ
	var songDetail SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &songDetail, nil
}
