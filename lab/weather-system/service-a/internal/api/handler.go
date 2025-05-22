package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"service-a/internal/config"
)

// NewRouter sets up the Gin router with all the routes
func NewRouter(cfg *config.Config, tracer trace.Tracer) *gin.Engine {
	router := gin.Default()
	router.POST("/weather", handleZipCodeInput(cfg, tracer))
	return router
}

// handleZipCodeInput handles the input ZIP code and forwards it to Service B
func handleZipCodeInput(cfg *config.Config, tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "handle_zipcode_input")
		defer span.End()

		var request struct {
			ZipCode string `json:"cep"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			span.RecordError(err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid zipcode"})
			return
		}

		// Validate zipcode format
		if !isValidZipCode(request.ZipCode) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid zipcode"})
			return
		}

		// Forward to Service B
		_, forwardSpan := tracer.Start(ctx, "forward_to_service_b")
		response, statusCode, err := forwardToServiceB(ctx, cfg.ServiceBURL, request.ZipCode)
		forwardSpan.End()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// Return response from Service B
		c.Data(statusCode, "application/json", response)
	}
}

func isValidZipCode(zipCode string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(zipCode)
}

func forwardToServiceB(ctx context.Context, serviceBURL, zipCode string) ([]byte, int, error) {
	requestBody, err := json.Marshal(map[string]string{
		"zipcode": zipCode,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", serviceBURL+"/weather", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Inject OpenTelemetry context into the request headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return responseBody, resp.StatusCode, nil
}
