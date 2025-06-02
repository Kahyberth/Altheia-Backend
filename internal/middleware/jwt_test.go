package middleware

import (
	"Altheia-Backend/pkg/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestMain(m *testing.M) {
	// Set up test environment variables
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-purposes")
	code := m.Run()
	os.Exit(code)
}

func TestJWTProtected(t *testing.T) {
	app := fiber.New()
	app.Use(JWTProtected())

	tests := []struct {
		name           string
		setupRequest   func(c *fiber.Ctx)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Missing Token",
			setupRequest: func(c *fiber.Ctx) {
				// No token set
			},
			expectedStatus: 401,
			expectedBody:   `{"error":"unauthorized: missing cookie"}`,
		},
		{
			name: "Invalid Token",
			setupRequest: func(c *fiber.Ctx) {
				c.Request().Header.Set("Cookie", "access_token=invalid.token.here")
			},
			expectedStatus: 401,
			expectedBody:   `{"error":"invalid token"}`,
		},
		{
			name: "Valid Token",
			setupRequest: func(c *fiber.Ctx) {
				token, _ := utils.GenerateJWT("test-user-123", 3600)
				c.Request().Header.Set("Cookie", "access_token="+token)
			},
			expectedStatus: 200,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request
			req := httptest.NewRequest("GET", "/", nil)
			_, _ = app.Test(req)

			// Setup the request context
			ctx := &fasthttp.RequestCtx{}
			c := app.AcquireCtx(ctx)
			tt.setupRequest(c)

			// Call the middleware
			err := JWTProtected()(c)

			// Check the response
			if tt.expectedStatus == 200 {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if c.Locals("user_id") != "test-user-123" {
					t.Errorf("Expected user_id 'test-user-123', got %v", c.Locals("user_id"))
				}
			} else {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if c.Response().StatusCode() != tt.expectedStatus {
					t.Errorf("Expected status %d, got %d", tt.expectedStatus, c.Response().StatusCode())
				}
				if string(c.Response().Body()) != tt.expectedBody {
					t.Errorf("Expected body %s, got %s", tt.expectedBody, string(c.Response().Body()))
				}
			}

			app.ReleaseCtx(c)
		})
	}
}

func TestJWTProtected_Integration(t *testing.T) {
	app := fiber.New()
	app.Use(JWTProtected())

	// Add a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("protected route")
	})

	// Test without token
	req := httptest.NewRequest("GET", "/test", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}

	// Test with invalid token
	req = httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: "invalid.token.here",
	})
	resp, _ = app.Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}

	// Test with valid token
	token, _ := utils.GenerateJWT("test-user-123", 3600)
	req = httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: token,
	})
	resp, _ = app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
