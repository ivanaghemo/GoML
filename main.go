package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	controller "main.go/controller"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.New()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Bienvenidos : Operación Fuego de Quasar!")
	})
	router.POST("/topsecret", controller.PostTopSecrets)

	router.Run(":" + port)

	log.Println("Listening..." + port)

}
