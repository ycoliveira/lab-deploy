package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock services for testing
func mockGetCityByZipCode(zipcode string) (string, error) {
	if zipcode == "01001000" {
		return "São Paulo", nil
	}
	if zipcode == "99999999" {
		return "", errors.New("CEP não encontrado")
	}
	return "", nil
}

func mockGetTemperatureByCity(city string) (float64, error) {
	if city == "São Paulo" {
		return 25.0, nil
	}
	return 0, errors.New("Cidade não encontrada")
}

func TestGetWeather(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/weather/:zipcode", func(c *gin.Context) {
		zipcode := c.Param("zipcode")

		if !isValidZipCode(zipcode) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
			return
		}

		city, err := mockGetCityByZipCode(zipcode)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "can not find zipcode"})
			return
		}

		tempC, err := mockGetTemperatureByCity(city)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get temperature"})
			return
		}

		response := gin.H{
			"temp_C": tempC,
			"temp_F": celsiusToFahrenheit(tempC),
			"temp_K": celsiusToKelvin(tempC),
		}
		c.JSON(http.StatusOK, response)
	})

	// Tests
	t.Run("Valid CEP", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/weather/01001000", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"temp_C": 25.0, "temp_F": 77.0, "temp_K": 298.15}`, w.Body.String())
	})

	t.Run("Invalid CEP Format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/weather/123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.JSONEq(t, `{"message": "invalid zipcode"}`, w.Body.String())
	})

	t.Run("CEP Not Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/weather/99999999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"message": "can not find zipcode"}`, w.Body.String())
	})
}

func TestIsValidZipCode(t *testing.T) {
	assert.True(t, isValidZipCode("01001000"))
	assert.False(t, isValidZipCode("0100100"))
	assert.False(t, isValidZipCode("abcdabcd"))
}

func TestCelsiusToFahrenheit(t *testing.T) {
	assert.Equal(t, 77.0, celsiusToFahrenheit(25.0))
	assert.Equal(t, 32.0, celsiusToFahrenheit(0.0))
	assert.Equal(t, 212.0, celsiusToFahrenheit(100.0))
}

func TestCelsiusToKelvin(t *testing.T) {
	assert.Equal(t, 298.15, celsiusToKelvin(25.0))
	assert.Equal(t, 273.15, celsiusToKelvin(0.0))
	assert.Equal(t, 373.15, celsiusToKelvin(100.0))
}
