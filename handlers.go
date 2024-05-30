package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handle(c *gin.Context) {
	zipcode := c.Param("zipcode")
	if !isValidZipCode(zipcode) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	location, err := getLocationByZipCode(zipcode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can not get location"})
		return
	}

	weather, err := getWeatherByLocation(location.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can not get weather"})
		return
	}

	response := getCurrentTemp(weather)
	c.JSON(http.StatusOK, response)
}
