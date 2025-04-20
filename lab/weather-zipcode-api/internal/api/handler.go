package api

import (
	"regexp"

	"github.com/gin-gonic/gin"

	"weather-zipcode-api/internal/config"
	"weather-zipcode-api/internal/location"
	"weather-zipcode-api/internal/models"
	"weather-zipcode-api/internal/weather"
)

// NewRouter sets up the Gin router with all the routes
func NewRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()
	router.POST("/weather", GetWeatherByZipCode(cfg))
	return router
}

// GetWeatherByZipCode handles requests to get weather data by zip code
func GetWeatherByZipCode(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			ZipCode string `json:"zipcode"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(422, gin.H{"error": "invalid zipcode"})
			return
		}

		zipCode := request.ZipCode

		// Validate zipcode format
		if !isValidZipCode(zipCode) {
			c.JSON(422, gin.H{"error": "invalid zipcode"})
			return
		}

		// Get location from zipcode
		address, err := location.GetAddress(zipCode)
		if err != nil {
			c.JSON(404, gin.H{"error": "can not find zipcode"})
			return
		}

		// Get weather data for the location - use the Localidade field directly
		tempC, err := weather.GetWeather(cfg.WeatherAPIKey, address.Localidade)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch weather data"})
			return
		}

		// Convert temperatures
		tempF := weather.ConvertToFahrenheit(tempC)
		tempK := weather.ConvertToKelvin(tempC)

		// Prepare response
		response := models.Weather{
			TempC: tempC,
			TempF: tempF,
			TempK: tempK,
		}

		c.JSON(200, response)
	}
}

func isValidZipCode(zipCode string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(zipCode)
}
