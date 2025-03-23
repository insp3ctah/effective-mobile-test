package song

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) AddSong(c *gin.Context) {
	var input Song
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	details, err := FetchSongDetails(input.Group, input.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	song := &Song{
		Group:       input.Group,
		Title:       input.Title,
		ReleaseDate: details.ReleaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}

	err = h.repo.Create(song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, song)
}

func (h *Handler) GetSongs(c *gin.Context) {
	group := c.Query("group")
	title := c.Query("title")
	releaseDate := c.Query("release_date")

	if group == "" && title == "" && releaseDate == "" {
		songs, err := h.repo.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, songs)
		return
	}

	filter := Song{
		Group:       group,
		Title:       title,
		ReleaseDate: releaseDate,
	}

	songs, err := h.repo.GetAllWithFilter(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}

func (h *Handler) DeleteSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.repo.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}

func (h *Handler) GetSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	song, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, song)
}

func (h *Handler) UpdateSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input Song
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	song, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	song.Group = input.Group
	song.Title = input.Title

	err = h.repo.Update(song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, song)
}

func (h *Handler) GetSongLyrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	song, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	lines := splitLyrics(song.Text)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(lines) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page number out of range"})
		return
	}
	if end > len(lines) {
		end = len(lines)
	}

	c.JSON(http.StatusOK, gin.H{
		"lyrics": lines[start:end],
		"page":   page,
		"size":   pageSize,
		"total":  len(lines),
	})
}

func splitLyrics(text string) []string {
	return strings.Split(text, "\n")
}
