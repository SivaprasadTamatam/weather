package main

import (
	"log"
	"net/http"

	"./weather"
)

// main is the entry point of the application.
// It sets up a simple HTTP server to handle incoming requests.
func main() {
	// Register the WeatherHandler function to handle requests to the "/weather" endpoint.
	// This is achieved using the built-in http package's HandleFunc method, which associates a handler function with a specific URL pattern.
	// For simplicity, we are using the basic capabilities of the standard http package instead of more advanced frameworks like GIN or MUX.
	http.HandleFunc("/weather", weather.WeatherHandler)

	// Start the HTTP server and listen for incoming requests on port 8080.
	// The ListenAndServe function is a blocking call, so the program will continue to run and serve requests until it is terminated.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
