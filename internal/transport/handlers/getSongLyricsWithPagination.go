package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// API URL (его можно вынести в .env)
const externalAPI = "http://external-api/info"

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// GetSongLyricsWithPagination получает текст песни с API и разбивает его на куплеты
func (h *SongHandler) GetSongLyricsWithPagination(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")

	log.Printf("[INFO] Received request for group: %s, song: %s", group, song)

	if group == "" || song == "" {
		log.Println("[ERROR] Missing required query parameters: group, song")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameters: group, song"})
		return
	}

	// Запрашиваем данные из внешнего API
	apiURL := fmt.Sprintf("%s?group=%s&song=%s", externalAPI, group, song)
	log.Printf("[DEBUG] Fetching song details from external API: %s", apiURL)

	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Failed to fetch song info from API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song info"})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var songDetail SongDetail
	if err := json.Unmarshal(body, &songDetail); err != nil {
		log.Printf("[ERROR] Error parsing response from API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing response from API"})
		return
	}

	log.Printf("[INFO] Successfully fetched song details for %s - %s", group, song)

	// Разбиваем текст на куплеты (по двойному переносу строки)
	verses := strings.Split(songDetail.Text, "\n\n")

	// Пагинация
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "1")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)

	if pageNum < 1 {
		pageNum = 1
	}
	if limitNum < 1 {
		limitNum = 1
	}

	start := (pageNum - 1) * limitNum
	end := start + limitNum

	if start >= len(verses) {
		log.Printf("[DEBUG] No verses available for page: %d", pageNum)
		c.JSON(http.StatusOK, gin.H{"verses": []string{}})
		return
	}
	if end > len(verses) {
		end = len(verses)
	}

	log.Printf("[INFO] Returning verses %d to %d out of %d total verses", start, end, len(verses))

	// Возвращаем результат
	c.JSON(http.StatusOK, gin.H{
		"group":       group,
		"song":        song,
		"releaseDate": songDetail.ReleaseDate,
		"link":        songDetail.Link,
		"verses":      verses[start:end],
		"page":        pageNum,
		"limit":       limitNum,
		"total":       len(verses),
	})
}
