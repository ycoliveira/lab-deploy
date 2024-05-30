package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

const viaCEPURL = "https://viacep.com.br/ws/%s/json/"
const weatherAPIKey = "ddaeea3388be423f813214917243005"
const weatherAPIURL = "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no"

type Location struct {
	Error    bool   `json:"erro"`
	Cep      string `json:"cep"`
	Location string `json:"localidade"`
}

type WeatherApiResp struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current Current `json:"current"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

type ResponseDto struct {
	TempC float64 `json:"temp_C"`
	TempK float64 `json:"temp_K"`
	TempF float64 `json:"temp_F"`
}

func isValidZipCode(zipcode string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(zipcode)
}

func getLocationByZipCode(zipcode string) (Location, error) {
	var location Location
	resp, err := http.Get(fmt.Sprintf(viaCEPURL, zipcode))
	if err != nil {
		log.Printf("Error fetching city: %v", err)
		return location, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("viaCEP responded with status code: %d", resp.StatusCode)
		return location, errors.New("failed to get city")
	}

	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		log.Printf("Error decoding viaCEP response: %v", err)
		return location, err
	}

	if location.Error {
		log.Printf("Error: invalid zipcode")
		return location, errors.New("can not find zipcode")
	}

	log.Printf("City found: %s", location.Location)
	return location, nil
}

func getWeatherByLocation(location string) (WeatherApiResp, error) {
	var weather WeatherApiResp
	cityEncoded := url.QueryEscape(location)
	url := fmt.Sprintf(weatherAPIURL, weatherAPIKey, cityEncoded)
	log.Printf("Requesting WeatherAPI with URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching temperature: %v", err)
		return weather, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("WeatherAPI responded with status code: %d", resp.StatusCode)
		return weather, fmt.Errorf("WeatherAPI responded with status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		log.Printf("Error decoding WeatherAPI response: %v", err)
		return weather, err
	}

	return weather, nil
}

func celsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func celsiusToKelvin(c float64) float64 {
	return c + 273.15
}

func getCurrentTemp(weather WeatherApiResp) ResponseDto {
	return ResponseDto{
		TempC: weather.Current.TempC,
		TempF: weather.Current.TempF,
		TempK: weather.Current.TempC + 273.15,
	}
}
