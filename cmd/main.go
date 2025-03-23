package main

import (
	"effective-mobile-test/config"
	"effective-mobile-test/internal/song"
	"effective-mobile-test/pkg"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func init() {
	pkg.InitDB()
	config.LoadConfig()
}

// @title Music API
// @version 1.0
// @description REST API for creating, reading, updating, and deleting songs
// @BasePath /

// @host localhost:8080
// @schemes http

func main() {
	repo := song.NewRepository()
	handler := song.NewHandler(repo)

	r := gin.Default()
	r.POST("/songs", handler.AddSong)
	r.GET("/songs", handler.GetSongs)
	r.DELETE("/songs/:id", handler.DeleteSong)
	r.GET("/songs/:id/lyrics", handler.GetSongLyrics)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("could not run server: %s", err.Error())
	}
}
