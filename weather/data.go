package weather

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const API_KEY = "346f820b8b57367bb052d099256c939d"

// getWeather is a function that retrieves weather data from the OpenWeatherMap API based on the provided latitude and longitude.
// It constructs the API URL using the latitude, longitude, and API key, and sends an HTTP GET request to fetch the data.
// If the HTTP request fails, it logs the error and returns nil and the error.
// If the JSON response from the API cannot be decoded, it logs the error and returns nil and the error.
// It then extracts relevant weather information such as description, temperature, visibility, wind speed, wind direction, cloud coverage, sunrise, and sunset from the JSON data.
// Finally, it constructs a WeatherData struct with the extracted information and returns it along with a nil error.
func getWeather(lat, lon float64) (*WeatherData, error) {
	// Construct the API URL reference https://openweathermap.org/current - API call section
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%.6f&lon=%.6f&appid=%s&units=metric", lat, lon, API_KEY)

	// Send HTTP GET request to the API
	response, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	// Decode the JSON response
	var data map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		return nil, err
	}

	// Extract weather information from the JSON data
	weatherDescription, temperature := extractWeatherInfo(data)
	visibility := extractVisibility(data)
	windSpeed, windDirection := extractWindInfo(data)
	cloudCoverage := extractCloudCoverage(data)
	sunrise, sunset := extractSunriseSunset(data)

	// Classify weather type based on temperature
	weatherType := classifyWeather(temperature)

	// Construct WeatherData struct and return
	return &WeatherData{
		WeatherDescription: weatherDescription,
		Temperature:        fmt.Sprintf("%v Celsius", temperature),
		WeatherType:        weatherType,
		Visibility:         visibility,
		WindSpeed:          windSpeed,
		WindDirection:      windDirection,
		CloudCoverage:      cloudCoverage,
		Sunrise:            sunrise,
		Sunset:             sunset,
	}, nil
}

// extractWeatherInfo is a helper function that extracts weather description and temperature from the JSON data.
func extractWeatherInfo(data map[string]interface{}) (string, float64) {
	// Extract weather description from the 'weather' field
	weatherArray := data["weather"].([]interface{})
	weatherDescription := weatherArray[0].(map[string]interface{})["description"].(string)

	// Extract temperature from the 'main' field
	temperature := data["main"].(map[string]interface{})["temp"].(float64)

	return weatherDescription, temperature
}

// extractVisibility is a helper function that extracts visibility from the JSON data.
func extractVisibility(data map[string]interface{}) string {
	// Extract visibility from the 'visibility' field and convert to kilometers
	visibility := int(data["visibility"].(float64)) / 1000
	return fmt.Sprintf("%v KM", visibility)
}

// extractWindInfo is a helper function that extracts wind speed and direction from the JSON data.
func extractWindInfo(data map[string]interface{}) (string, string) {
	// Extract wind speed and direction from the 'wind' field
	windData := data["wind"].(map[string]interface{})
	windSpeed := windData["speed"].(float64)
	windDirection := int(windData["deg"].(float64))
	return fmt.Sprintf("%v meter/sec", windSpeed), fmt.Sprintf("%v degrees", windDirection)
}

// extractCloudCoverage is a helper function that extracts cloud coverage from the JSON data.
func extractCloudCoverage(data map[string]interface{}) string {
	// Extract cloud coverage from the 'clouds' field
	cloudData := data["clouds"].(map[string]interface{})
	cloudCoverage := int(cloudData["all"].(float64))
	return fmt.Sprintf("%v percentage", cloudCoverage)
}

// extractSunriseSunset is a helper function that extracts sunrise and sunset times from the JSON data.
func extractSunriseSunset(data map[string]interface{}) (time.Time, time.Time) {
	// Extract sunrise and sunset times from the 'sys' field
	sunriseUnix := int64(data["sys"].(map[string]interface{})["sunrise"].(float64))
	sunsetUnix := int64(data["sys"].(map[string]interface{})["sunset"].(float64))
	sunrise := time.Unix(sunriseUnix, 0)
	sunset := time.Unix(sunsetUnix, 0)
	return sunrise, sunset
}

// classifyWeather is a helper function that classifies the weather type based on temperature.
func classifyWeather(temperature float64) string {
	// Classify weather type based on temperature ranges
	if temperature <= 10 {
		return "cold"
	} else if temperature <= 25 {
		return "moderate"
	}
	return "hot"
}
