package weather

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// WeatherData represents the structure of weather data obtained from the OpenWeatherMap API.
// It is constructed based on the JSON response format documented at https://openweathermap.org/current.
type WeatherData struct {
	WeatherDescription string    `json:"weather_condition"` // Description of the weather condition
	Temperature        string    `json:"temperature"`       // Temperature in Celsius
	WeatherType        string    `json:"weather_type"`      // Type of weather condition (e.g., cold, moderate, hot)
	Visibility         string    `json:"visibility"`        // Visibility in kilometers
	WindSpeed          string    `json:"wind_speed"`        // Wind speed in meters per second
	WindDirection      string    `json:"wind_direction"`    // Wind direction in degrees
	CloudCoverage      string    `json:"cloud_coverage"`    // Cloud coverage in percentage
	Sunrise            time.Time `json:"sunrise"`           // Time of sunrise
	Sunset             time.Time `json:"sunset"`            // Time of sunset
}

// WeatherHandler is an HTTP handler function that processes incoming HTTP requests to fetch weather data.
// It expects latitude and longitude parameters in the request URL query string.
// If the latitude or longitude parameters are missing or invalid, it responds with a Bad Request status code (400).
// It then calls the getWeatherWithContext function to retrieve weather data based on the provided latitude and longitude.
// If there is an error during the weather data retrieval process, it responds with an Internal Server Error status code (500).
// Otherwise, it encodes the retrieved weather data into JSON format and writes it to the response writer.
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	// Parse latitude and longitude from the request URL query parameters
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call getWeatherWithContext function with the created context
	weatherData, err := getWeatherWithContext(ctx, lat, lon)
	if err != nil {
		// Handle error if any occurred during weather data retrieval
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	// Encode weather data into JSON format and write it to the response writer
	json.NewEncoder(w).Encode(weatherData)
}

// getWeatherWithContext retrieves weather data with a deadline context
func getWeatherWithContext(ctx context.Context, lat, lon float64) (*WeatherData, error) {
	// Create channels to communicate results and errors
	ch := make(chan *WeatherData, 1)
	errCh := make(chan error, 1)

	// Execute getWeather function asynchronously
	go func() {
		weatherData, err := getWeather(lat, lon)
		if err != nil {
			// Send error to the error channel if any occurred
			errCh <- err
			return
		}
		// Send weather data to the result channel
		ch <- weatherData
	}()

	// Select block to wait for results or errors
	select {
	case <-ctx.Done():
		// Return error if context deadline is reached
		return nil, ctx.Err()
	case err := <-errCh:
		// Return error if any occurred during weather data retrieval
		return nil, err
	case weatherData := <-ch:
		// Return weather data if retrieved successfully
		return weatherData, nil
	}
}
