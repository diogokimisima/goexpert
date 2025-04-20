# Weather Zipcode API

This project is a simple API built in Go that allows users to retrieve current weather information based on a provided ZIP code. The API utilizes the ViaCEP service to fetch location data and the WeatherAPI to obtain weather information.

## Features

- Accepts a valid 8-digit ZIP code.
- Returns the current temperature in Celsius, Fahrenheit, and Kelvin.
- Handles errors for invalid or non-existent ZIP codes.

## Project Structure

```
weather-zipcode-api
├── cmd
│   └── server
│       └── main.go          # Entry point of the application
├── internal
│   ├── api
│   │   └── handler.go       # HTTP handlers for the API
│   ├── config
│   │   └── config.go        # Configuration management
│   ├── location
│   │   └── viacep.go        # Interaction with ViaCEP API
│   ├── models
│   │   └── weather.go       # Data structures for weather information
│   └── weather
│       └── weatherapi.go    # Interaction with WeatherAPI
├── tests
│   └── api_test.go          # Automated tests for API handlers
├── Dockerfile                # Docker image build instructions
├── docker-compose.yml        # Service definitions for Docker Compose
├── go.mod                   # Module definition and dependencies
└── go.sum                   # Dependency checksums
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone https://github.com/yourusername/weather-zipcode-api.git
   cd weather-zipcode-api
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Configure environment variables:**
   Set up your API keys and any other necessary configuration in the `config.go` file.

4. **Run the application:**
   ```
   go run cmd/server/main.go
   ```

5. **Access the API:**
   The API will be available at `http://localhost:8080`. You can test it by sending requests to the appropriate endpoints.

## Usage

To get the current weather for a ZIP code, send a GET request to:
```
GET /weather?zipcode={ZIP_CODE}
```

### Example Response
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

## Error Handling

- **422 Unprocessable Entity**: Returned when the ZIP code format is invalid.
- **404 Not Found**: Returned when the ZIP code cannot be found.

## Deployment

This application can be deployed on Google Cloud Run. Follow the [Google Cloud documentation](https://cloud.google.com/run/docs/deploying) for instructions on deploying a Docker container.

You can access the live API at: [https://weather-zipcode-api-999076993121.us-central1.run.app](https://weather-zipcode-api-999076993121.us-central1.run.app)

To use the API, send a POST request to the `/weather` endpoint with a JSON body containing the zipcode. For example:

This project is licensed under the MIT License. See the LICENSE file for more details.