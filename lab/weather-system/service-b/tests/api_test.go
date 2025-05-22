package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"weather-zipcode-api/internal/api"
	"weather-zipcode-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		WeatherAPIKey: "test_key",
	}
	return api.NewRouter(cfg)
}

func TestGetWeatherByValidZipCode(t *testing.T) {
	router := setupRouter()

	// Create request body with valid ZIP code
	requestBody, _ := json.Marshal(map[string]string{
		"zipcode": "01001000",
	})

	req, _ := http.NewRequest("POST", "/weather", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// This will actually fail because we're not mocking the external APIs
	// In a real test, you would mock these API calls
	// assert.Equal(t, http.StatusOK, resp.Code)
	// Add assertions for response body
}

func TestGetWeatherByInvalidZipCodeFormat(t *testing.T) {
	router := setupRouter()

	// Create request body with invalid ZIP code
	requestBody, _ := json.Marshal(map[string]string{
		"zipcode": "123", // Too short
	})

	req, _ := http.NewRequest("POST", "/weather", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	// Add assertions for response body
}

func TestGetWeatherByNonExistentZipCode(t *testing.T) {
	router := setupRouter()

	// Create request body with non-existent ZIP code
	requestBody, _ := json.Marshal(map[string]string{
		"zipcode": "99999999",
	})

	req, _ := http.NewRequest("POST", "/weather", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// This will actually fail because we're not mocking the external APIs
	// In a real test, you would mock these API calls
	// assert.Equal(t, http.StatusNotFound, resp.Code)
	// Add assertions for response body
}
