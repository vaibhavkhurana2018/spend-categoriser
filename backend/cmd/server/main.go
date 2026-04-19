package main

import (
	"gofr.dev/pkg/gofr"

	"github.com/spend-categoriser/backend/internal/handler"
)

func main() {
	app := gofr.New()

	app.GET("/api/health", handler.Health)
	app.POST("/api/parse", handler.ParsePDF)
	app.GET("/api/categories", handler.GetCategories)

	app.Run()
}
