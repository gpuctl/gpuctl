package webapi_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gpuctl/gpuctl/internal/authentication"
	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"github.com/gpuctl/gpuctl/internal/webapi"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	mockLogger := slog.Default()

	mockRequest := httptest.NewRequest(http.MethodPost, "/api/auth", nil)

	auth := webapi.ConfigFileAuthenticator{
		Username:      "joe",
		Password:      "mama",
		CurrentTokens: make(map[authentication.AuthToken]bool),
	}
	creds := webapi.APIAuthCredientals{Username: "joe", Password: "mama"}

	mockDB := database.InMemory()

	api := &webapi.Api{DB: mockDB}

	response, err := api.Authenticate(&auth, creds, mockRequest, mockLogger)

	found := false
	for _, cookie := range response.Cookies {
		if cookie.Name == authentication.TokenCookieName {
			found = true
			break
		}
	}
	assert.True(t, found, "Cookies contain auth cookie")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestLogOut(t *testing.T) {
	mockLogger := slog.Default()

	auth := webapi.ConfigFileAuthenticator{
		Username:      "joe",
		Password:      "mama",
		CurrentTokens: make(map[authentication.AuthToken]bool),
	}
	creds := webapi.APIAuthCredientals{Username: "joe", Password: "mama"}

	mockDB := database.InMemory()

	api := &webapi.Api{DB: mockDB}

	token, err := auth.CreateToken(creds)
	assert.NoError(t, err, "No error in creating auth token")

	// Make new response to revoke the token
	revokeRequest := httptest.NewRequest(http.MethodGet, "/api/admin/logout", nil)
	revokeRequest.AddCookie(&http.Cookie{Name: authentication.TokenCookieName, Value: token})

	response, err := api.LogOut(&auth, revokeRequest, mockLogger)
	assert.NoError(t, err, "No error in logging-out")
	assert.Equal(t, http.StatusOK, response.Status)

	unauthenticatedRequest := httptest.NewRequest(http.MethodGet, "/api/admin/confirm", nil)
	unauthenticatedRequest.AddCookie(&http.Cookie{Name: authentication.TokenCookieName, Value: token})
	resp, err := api.ConfirmAdmin(&auth, unauthenticatedRequest, mockLogger)
	assert.NoError(t, err, "No error in unauthenticated request")
	assert.Equal(t, http.StatusUnauthorized, resp.Status)
}

func TestAllStatistics(t *testing.T) {
	// Mock GPUInfo and GPUStatSample data
	gpuInfos := []uplink.GPUInfo{
		{
			Uuid:          "uuid-123",
			Name:          "GeForce GTX 1080",
			Brand:         "NVIDIA",
			DriverVersion: "441.66",
			MemoryTotal:   8192, // 8GB
		},
	}

	stats := []uplink.GPUStatSample{
		{
			Uuid:              "uuid-123",
			MemoryUtilisation: 64.5,
			GPUUtilisation:    75.2,
			MemoryUsed:        4096, // 4GB
			FanSpeed:          55.0,
			Temp:              70.0, // 70°C
			MemoryTemp:        68.0, // 68°C
			GraphicsVoltage:   1.05,
			PowerDraw:         150.0,                  // 150 Watts
			GraphicsClock:     1750.0,                 // 1750 Mhz
			MaxGraphicsClock:  1800.0,                 // 1800 Mhz
			MemoryClock:       5000.0,                 // 5000 Mhz
			MaxMemoryClock:    5100.0,                 // 5100 Mhz
			RunningProcesses:  []uplink.GPUProcInfo{}, // Assume GPUProcInfo is defined elsewhere
		},
	}

	tests := []struct {
		name          string
		mockData      []uplink.GpuStatsUpload
		mockError     error
		expectedError bool
	}{
		{
			name: "successful data retrieval",
			mockData: []uplink.GpuStatsUpload{
				{
					Hostname: "host1",
					Stats:    stats,
					GPUInfos: gpuInfos,
				},
			},
			expectedError: false,
		},
		{
			name:          "empty data",
			mockData:      []uplink.GpuStatsUpload{},
			expectedError: false,
		},
	}

	logger := slog.Default()
	mockDB := database.InMemory()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			api := webapi.Api{DB: mockDB}
			req, _ := http.NewRequest("GET", "/url", nil)

			_, err := api.AllStatistics(req, logger)

			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}

			// ! Somebody will need to write the tests for what we expect as output data
		})
	}
}

func TestZipStats(t *testing.T) {
	host := "testHost"
	info := uplink.GPUInfo{
		Uuid:          "uuid-123",
		Name:          "GeForce GTX 1080",
		Brand:         "NVIDIA",
		DriverVersion: "441.66",
		MemoryTotal:   8192,
	}
	stat := uplink.GPUStatSample{
		Uuid:              "uuid-123",
		MemoryUtilisation: 64.5,
		GPUUtilisation:    75.2,
		MemoryUsed:        4096,
		FanSpeed:          55.0,
		Temp:              70.0,
		MemoryTemp:        68.0,
		GraphicsVoltage:   1.05,
		PowerDraw:         150.0,
		GraphicsClock:     1750.0,
		MaxGraphicsClock:  1800.0,
		MemoryClock:       5000.0,
		MaxMemoryClock:    5100.0,
	}

	expected := broadcast.OldGPUStatSample{
		Hostname:          host,
		Uuid:              info.Uuid,
		Name:              info.Name,
		Brand:             info.Brand,
		DriverVersion:     info.DriverVersion,
		MemoryTotal:       info.MemoryTotal,
		MemoryUtilisation: stat.MemoryUtilisation,
		GPUUtilisation:    stat.GPUUtilisation,
		MemoryUsed:        stat.MemoryUsed,
		FanSpeed:          stat.FanSpeed,
		Temp:              stat.Temp,
		MemoryTemp:        stat.MemoryTemp,
		GraphicsVoltage:   stat.GraphicsVoltage,
		PowerDraw:         stat.PowerDraw,
		GraphicsClock:     stat.GraphicsClock,
		MaxGraphicsClock:  stat.MaxGraphicsClock,
		MemoryClock:       stat.MemoryClock,
		MaxMemoryClock:    stat.MaxMemoryClock,
	}

	result := webapi.ZipStats(host, info, stat)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ZipStats() = %v, want %v", result, expected)
	}
}

func TestServerEndpoints(t *testing.T) {
	mockDB := database.InMemory()

	auth := webapi.ConfigFileAuthenticator{
		Username:      "joe",
		Password:      "mama",
		CurrentTokens: map[authentication.AuthToken]bool{"example_token": true},
	}

	server := webapi.NewServer(mockDB, &auth, webapi.OnboardConf{})

	tests := []struct {
		name           string
		method         string
		endpoint       string
		body           []byte
		expectedStatus int
		headers        map[string]string
	}{
		{
			name:           "Test Stats All",
			method:         http.MethodGet,
			endpoint:       "/api/stats/all",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test Offline Machines",
			method:         http.MethodGet,
			endpoint:       "/api/stats/offline",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test Confirm Authentication Fails",
			method:         http.MethodGet,
			endpoint:       "/api/admin/confirm",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Test Authentication Fails",
			method:         http.MethodPost,
			endpoint:       "/api/admin/auth",
			body:           []byte(`{"username":"mama","password":"joe"}`),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Test Authentication",
			method:         http.MethodPost,
			endpoint:       "/api/admin/auth",
			body:           []byte(`{"username":"joe","password":"mama"}`),
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "Test Confirm Authentication Succeeds",
			method:         http.MethodGet,
			endpoint:       "/api/admin/confirm",
			expectedStatus: http.StatusOK,
			headers:        map[string]string{"Cookie": "token=example_token"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, tc.endpoint, bytes.NewBuffer(tc.body))
			for k, v := range tc.headers {
				request.Header.Add(k, v)
			}
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, request)

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("%s: expected status code %d, got %d", tc.name, tc.expectedStatus, status)
			}
		})
	}
}
