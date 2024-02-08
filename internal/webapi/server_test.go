package webapi_test

// import (
// 	"log/slog"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gpuctl/gpuctl/internal/authentication"
// 	"github.com/gpuctl/gpuctl/internal/database"
// 	"github.com/gpuctl/gpuctl/internal/webapi"
// )

// func TestAuthenticate(t *testing.T) {
// 	// Create a mock logger for testing
// 	mockLogger := slog.Default()

// 	// Create a mock request for testing
// 	mockRequest := httptest.NewRequest(http.MethodPost, "/api/auth", nil)

// 	// Create a ConfigFileAuthenticator instance for testing
// 	auth := webapi.ConfigFileAuthenticator{
// 		Username:      "joe",
// 		Password:      "mama",
// 		CurrentTokens: make(map[authentication.AuthToken]bool),
// 	}

// 	// Create APIAuthCredientals for testing
// 	creds := webapi.APIAuthCredientals{Username: "joe", Password: "mama"}

// 	// Create a mock database instance (if needed for testing)
// 	mockDB := database.InMemory()

// 	// Create an instance of the API struct with the mock database
// 	api := &webapi.Api{DB: mockDB}

// 	// Call the Authenticate function with the mock authenticator, credentials, mock request, and mock logger
// 	token, err := api.Authenticate(auth, creds, mockRequest, mockLogger)

// 	// Check if the returned token is valid
// 	if token == "" {
// 		t.Error("Expected a valid authentication token, got empty string")
// 	}

// 	// Check if there is no error returned
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// }
