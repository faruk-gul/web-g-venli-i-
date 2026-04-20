package main

import (
	"log"
	"os"

	"github.com/faruk/secscan/backend/internal/api"
	"github.com/faruk/secscan/backend/internal/scanner"
	"github.com/gin-gonic/gin"
)

func main() {
	service := scanner.NewService(scanner.NewRegistry())
	router := gin.Default()
	api.RegisterRoutes(router, service)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("secscan backend listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

