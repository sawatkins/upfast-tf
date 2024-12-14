package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sawatkins/upfast-tf/database"
)

// Test NotFound returns 404 and return html
func TestNotFound(t *testing.T) {
	engine := html.New("../templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(NotFound)

	req := httptest.NewRequest(http.MethodGet, "/non-existent-route", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check if the status code is 404
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", resp.StatusCode)
	}

	// Check that the content type response is html
	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected content type 'text/html; charset=utf-8', got '%s'", contentType)
	}
}

// Test Index route existance, status code, and content-type
func TestIndex(t *testing.T) {
	database.InitDB(":memory:")
	database.InitPlayerSessionTable()

	engine := html.New("../templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/", Index)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Check that the content type response is html
	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected content type 'text/html; charset=utf-8', got '%s'", contentType)
	}
}

// Test About route existance, status code, and content-type
func TestAbout(t *testing.T) {
	engine := html.New("../templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/about", About)

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Check that the content type response is html
	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected content type 'text/html; charset=utf-8', got '%s'", contentType)
	}
}

func TestGetServerIPs(t *testing.T) {
	// Initialize test database
	database.InitDB(":memory:")
	database.InitServerTable()

	engine := html.New("../templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/api/servers", GetServerIPs)

	// Test 1: Empty server list
	req := httptest.NewRequest(http.MethodGet, "/api/servers", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check status code is 200 even with empty list
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Test 2: With server IPs in database
	// Insert test data using ExecuteSQL
	insertSQL := `
		INSERT INTO servers (
			instance_id, public_ip, public_dns, name, server_hostname, map, players, max_players
		) VALUES 
		('i-1234567890', '192.168.1.1', 'ec2-1.compute.amazonaws.com', 'Server1', 'TF2 Server 1', 'cp_badlands', 0, 24),
		('i-0987654321', '192.168.1.2', 'ec2-2.compute.amazonaws.com', 'Server2', 'TF2 Server 2', 'cp_process', 0, 24)
	`
	database.ExecuteSQL(insertSQL)

	req = httptest.NewRequest(http.MethodGet, "/api/servers", nil)
	resp, err = app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Check content type
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got '%s'", contentType)
	}

	// Test 3: Database error scenario
	// Close the database to simulate an error
	database.Close()

	req = httptest.NewRequest(http.MethodGet, "/api/servers", nil)
	resp, err = app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check error status code
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", resp.StatusCode)
	}
}

GetServerInfo