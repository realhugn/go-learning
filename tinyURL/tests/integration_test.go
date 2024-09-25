package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"tinyURL/config"
	"tinyURL/database"
	"tinyURL/internal/api"
	"tinyURL/internal/handlers"
	"tinyURL/internal/repository"
	"tinyURL/internal/service"

	"github.com/gin-gonic/gin"
)

func TestTinyURLIntegration(t *testing.T) {
	// Setup
	cfg := &config.Config{
		DatabaseURL: "postgres://test:123@localhost:5432/tinyurl_test",
		ServerPort:  "8082",
	}

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	redisClient, _ := database.NewRedisCache(cfg)

	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo, redisClient)
	urlHandler := handlers.NewURLHandler(urlService)

	router := gin.Default()
	api.SetupRoutes(router, urlHandler)

	// Test cases
	tests := []struct {
		name           string
		method         string
		path           string
		body           map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Create Short URL",
			method:         "POST",
			path:           "/shorten",
			body:           map[string]string{"long_url": "https://www.google.com/"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"short_url":"http://localhost:8082/[a-zA-Z0-9]{8}"}`,
		},
		{
			name:           "Get Original URL",
			method:         "GET",
			path:           "/{short_url}",
			expectedStatus: http.StatusFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tc.method == "POST" {
				bodyJSON, _ := json.Marshal(tc.body)
				req, err = http.NewRequest(tc.method, tc.path, bytes.NewBuffer(bodyJSON))
			} else {
				req, err = http.NewRequest(tc.method, tc.path, nil)
			}

			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			if tc.method == "POST" {
				var response map[string]string
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Could not parse response body: %v", err)
				}
				short_url, ok := response["short_url"]
				if !ok {
					t.Errorf("Response does not contain short_url")
				}
				short_url = short_url[22:]
				tests[1].path = "/" + short_url
			}

			if tc.method == "GET" {
				location := rr.Header().Get("Location")

				if location != "https://www.google.com/" {
					t.Errorf("Expected Location header to be %s, got %s", "https://www.google.com/", location)
				}
			}
		})
	}
}
