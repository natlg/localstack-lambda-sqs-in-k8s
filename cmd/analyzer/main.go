package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var msgReceived uint

func main() {
	r := gin.Default()

	r.GET("/analyze/msg", func(c *gin.Context) {
		msgReceived++
		log.Printf("analyzed: %v", msgReceived)
		c.JSON(http.StatusOK, gin.H{"analyzed": msgReceived})
	})

	r.GET("/analyze/stats", func(c *gin.Context) {
		log.Printf("stats, msgReceived %v", msgReceived)
		c.JSON(http.StatusOK, gin.H{"analyzerStats": msgReceived})
	})
	if err := r.Run(":8081"); err != nil {
		return
	}
}
