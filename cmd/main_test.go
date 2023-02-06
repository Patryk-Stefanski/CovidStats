package main

import (
	"CovidStats2/client"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestIndexHandler(t *testing.T) {
	myClient := &http.Client{Timeout: 10 * time.Second}

	err := os.Setenv("ENVIRONMENT", "test")
	if err != nil {
		log.Fatalf("err %v", err)
	}

	err = godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: Failed to load .env file")
	}

	apiKey := os.Getenv("COVID_STATS_API_KEY")
	statsApi := client.NewClient(myClient, apiKey)

	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
		return
	}

	response := httptest.NewRecorder()
	handler := indexHandler(statsApi)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.Code)
		return
	}

	expectedResponse := "<option  value=Ireland>Ireland</option>"
	if !strings.Contains(response.Body.String(), expectedResponse) {
		t.Errorf("Expected response to contain %q, but got %q", expectedResponse, response.Body.String())
		return
	}
}

func TestHandleLive(t *testing.T) {
	myClient := &http.Client{Timeout: 10 * time.Second}

	err := os.Setenv("ENVIRONMENT", "test")
	if err != nil {
		log.Fatalf("err %v", err)
	}

	err = godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: Failed to load .env file")
	}

	apiKey := os.Getenv("COVID_STATS_API_KEY")
	statsApi := client.NewClient(myClient, apiKey)

	url := "/searchLive?country=Ireland"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
		return
	}

	response := httptest.NewRecorder()
	handler := HandleLive(statsApi)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.Code)
		return
	}

	expectedResponse := "Ireland"
	if !strings.Contains(response.Body.String(), expectedResponse) {
		t.Errorf("Expected response to contain %q, but got %q", expectedResponse, response.Body.String())
		return
	}
}

func TestHandleHistorical(t *testing.T) {
	myClient := &http.Client{Timeout: 10 * time.Second}

	err := os.Setenv("ENVIRONMENT", "test")
	if err != nil {
		log.Fatalf("err %v", err)
	}

	err = godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: Failed to load .env file")
	}

	apiKey := os.Getenv("COVID_STATS_API_KEY")
	statsApi := client.NewClient(myClient, apiKey)

	url := "/searchHistorical?country=USA&date=2023-02-03"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
		return
	}

	response := httptest.NewRecorder()
	handler := HandleHistorical(statsApi)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.Code)
	}

	expectedResponse := "2023-02-03"
	if !strings.Contains(response.Body.String(), expectedResponse) {
		t.Errorf("Expected response to contain %q, but got %q", expectedResponse, response.Body.String())
		return
	}
}
