package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/weather/:zipcode", Handle)
	log.Fatal(router.Run(":8080"))
}
