package api

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"weather-zipcode-api/internal/config"
	"weather-zipcode-api/internal/location"
	"weather-zipcode-api/internal/weather"
)

// NewRouter sets up the Gin router with all the routes
func NewRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(otelMiddleware())
	router.POST("/weather", GetWeatherByZipCode(cfg))
	return router
}

// otelMiddleware extracts the OTel context from HTTP headers
func otelMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GetWeatherByZipCode handles requests to get weather data by zip code
func GetWeatherByZipCode(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tracer
		tracer := otel.Tracer("service-b-tracer")
		ctx := c.Request.Context()

		// Create main span
		ctx, span := tracer.Start(ctx, "get_weather_by_zipcode")
		defer span.End()

		var request struct {
			ZipCode string `json:"zipcode"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			span.RecordError(err)
			c.JSON(422, gin.H{"error": "invalid zipcode"})
			return
		}

		zipCode := request.ZipCode

		// Validate zipcode format
		if !isValidZipCode(zipCode) {
			c.JSON(422, gin.H{"error": "invalid zipcode"})
			return
		}

		// Get location from zipcode with span
		ctx, locSpan := tracer.Start(ctx, "get_location_from_zipcode")
		address, err := location.GetAddress(zipCode)
		locSpan.End()

		if err != nil {
			c.JSON(404, gin.H{"error": "can not find zipcode"})
			return
		}

		// Get weather data for the location with span
		ctx, weatherSpan := tracer.Start(ctx, "get_weather_data")
		tempC, err := weather.GetWeather(cfg.WeatherAPIKey, address.Localidade)
		weatherSpan.End()

		if err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch weather data"})
			return
		}

		// Convert temperatures
		tempF := weather.ConvertToFahrenheit(tempC)
		tempK := weather.ConvertToKelvin(tempC)

		// Prepare response
		response := struct {
			City  string  `json:"city"`
			TempC float64 `json:"temp_C"`
			TempF float64 `json:"temp_F"`
			TempK float64 `json:"temp_K"`
		}{
			City:  address.Localidade,
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
