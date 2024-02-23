package webapi_test

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gpuctl/gpuctl/internal/authentication"
	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/tunnel"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"github.com/gpuctl/gpuctl/internal/webapi"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/uploadtest.pdf
var uploadPdfBytes []byte
var uploadPdfEnc = base64.StdEncoding.EncodeToString(uploadPdfBytes)

//go:embed testdata/more.txt
var uploadTxtBytes []byte
var uploadTxtEnc = base64.StdEncoding.EncodeToString(uploadTxtBytes)

var (
	joeAuth = webapi.ConfigFileAuthenticator{
		Username:      "joe",
		Password:      "mama",
		CurrentTokens: make(map[authentication.AuthToken]bool),
	}
	joeCreds = webapi.APIAuthCredientals{Username: "joe", Password: "mama"}
)

func TestAuthenticate(t *testing.T) {
	mockLogger := slog.Default()

	mockRequest := httptest.NewRequest(http.MethodPost, "/api/auth", nil)

	mockDB := database.InMemory()

	api := &webapi.Api{DB: mockDB}

	response, err := api.Authenticate(&joeAuth, joeCreds, mockRequest, mockLogger)

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
	mockDB := database.InMemory()

	api := &webapi.Api{DB: mockDB}

	token, err := joeAuth.CreateToken(joeCreds)
	assert.NoError(t, err, "No error in creating auth token")

	// Make new response to revoke the token
	revokeRequest := httptest.NewRequest(http.MethodGet, "/api/admin/logout", nil)
	revokeRequest.AddCookie(&http.Cookie{Name: authentication.TokenCookieName, Value: token})

	response, err := api.LogOut(&joeAuth, revokeRequest, mockLogger)
	assert.NoError(t, err, "No error in logging-out")
	assert.Equal(t, http.StatusOK, response.Status)

	unauthenticatedRequest := httptest.NewRequest(http.MethodGet, "/api/admin/confirm", nil)
	unauthenticatedRequest.AddCookie(&http.Cookie{Name: authentication.TokenCookieName, Value: token})
	resp, err := api.ConfirmAdmin(&joeAuth, unauthenticatedRequest, mockLogger)
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

func TestListFiles(t *testing.T) {
	mockDB := database.InMemory()
	mockLogger := slog.Default()
	api := &webapi.Api{DB: mockDB}
	hostname := "machine01"
	mockDB.UpdateLastSeen(hostname, 0)

	token, err := joeAuth.CreateToken(joeCreds)
	assert.NoError(t, err, "No error in creating auth token")
	cookie := http.Cookie{Name: authentication.TokenCookieName, Value: token}

	pdf1 := broadcast.AttachFile{
		Hostname:    hostname,
		Filename:    "file1",
		Mime:        "application/pdf",
		EncodedFile: uploadPdfEnc,
	}

	pdf2 := broadcast.AttachFile{
		Hostname:    hostname,
		Filename:    "file2",
		Mime:        "application/pdf",
		EncodedFile: uploadPdfEnc,
	}

	req := httptest.NewRequest(http.MethodPost, "/api/admin/attach_file", nil)
	req.AddCookie(&cookie)
	api.AttachFile(pdf1, req, mockLogger)
	resp, err := api.AttachFile(pdf2, req, mockLogger)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)

	req = httptest.NewRequest(http.MethodGet, "/api/admin/list_files?hostname="+hostname, nil)
	req.AddCookie(&cookie)

	listresp, err := api.ListFiles(req, mockLogger)
	assert.Equal(t, http.StatusOK, listresp.Status)
	list := listresp.Body
	assert.ElementsMatch(t, list, []string{pdf1.Filename, pdf2.Filename})
}

func TestRemovingFile(t *testing.T) {
	mockDB := database.InMemory()
	mockLogger := slog.Default()
	api := &webapi.Api{DB: mockDB}
	hostname := "machine09"

	token, err := joeAuth.CreateToken(joeCreds)
	assert.NoError(t, err, "No error in creating auth token")
	cookie := http.Cookie{Name: authentication.TokenCookieName, Value: token}

	mockDB.UpdateLastSeen(hostname, 0)

	pdf := broadcast.AttachFile{
		Hostname:    hostname,
		Filename:    "verycoolfile",
		Mime:        "application/pdf",
		EncodedFile: uploadPdfEnc,
	}

	// Add file
	req := httptest.NewRequest(http.MethodPost, "/api/admin/attach_file", nil)
	req.AddCookie(&cookie)
	_, err = api.AttachFile(pdf, req, mockLogger)
	assert.NoError(t, err)

	// Remove file
	remreq := httptest.NewRequest(http.MethodPost, "/api/admin/remove_file", nil)
	req.AddCookie(&cookie)
	res, err := api.RemoveFile(broadcast.RemoveFile{Hostname: hostname, Filename: pdf.Filename}, remreq, mockLogger)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Status)

	// Confirm from list
	listreq := httptest.NewRequest(http.MethodGet, "/api/admin/list_files?hostname="+hostname, nil)
	listreq.AddCookie(&cookie)
	listresp, err := api.ListFiles(listreq, mockLogger)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, listresp.Status)
	list := listresp.Body
	assert.Equal(t, list, []string{})
}

func TestAttachingFile(t *testing.T) {
	mockDB := database.InMemory()
	mockLogger := slog.Default()
	api := &webapi.Api{DB: mockDB}

	mockDB.UpdateLastSeen("testmachine", 0)

	token, err := joeAuth.CreateToken(joeCreds)
	assert.NoError(t, err, "No error in creating auth token")
	cookie := http.Cookie{Name: authentication.TokenCookieName, Value: token}
	payload := broadcast.AttachFile{
		Hostname:    "testmachine",
		Filename:    "testfile",
		Mime:        "application/pdf",
		EncodedFile: uploadPdfEnc,
	}

	// request for adding file
	req := httptest.NewRequest(http.MethodPost, "/api/admin/attach_file", nil)
	req.AddCookie(&cookie)
	resp, err := api.AttachFile(payload, req, mockLogger)
	assert.NoError(t, err, "No error in valid request to attach a file")
	if err != nil {
		return
	}
	assert.Equal(t, http.StatusOK, resp.Status)

	// request for getting file
	req = httptest.NewRequest(http.MethodGet, "/api/admin/get_file?hostname=testmachine&file=testfile", nil)
	req.AddCookie(&cookie)
	getresp, err := api.GetFile(req, mockLogger)
	assert.NoError(t, err, "No error in valid request to download file")
	if err != nil {
		return
	}
	assert.Equal(t, http.StatusOK, getresp.Status)
	// Compare bytes of the files
	assert.Equal(t, uploadPdfBytes, getresp.Body)
}

func TestServerEndpoints(t *testing.T) {
	mockDB := database.InMemory()

	auth := webapi.ConfigFileAuthenticator{
		Username:      "joe",
		Password:      "mama",
		CurrentTokens: map[authentication.AuthToken]bool{"example_token": true},
	}

	server := webapi.NewServer(mockDB, &auth, tunnel.Config{})

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
		{
			name:           "Test uploading file is authenticated",
			method:         http.MethodPost,
			endpoint:       "/api/admin/attach_file",
			expectedStatus: http.StatusUnauthorized,
			body:           []byte(`{"hostname":"n/a", "mime":"n/a", "file_enc":""}`),
			headers:        map[string]string{"Cookie": "token=wrongtoken"},
		},
		{
			name:           "Test downloading file is authenticated",
			method:         http.MethodGet,
			endpoint:       "/api/admin/get_file?hostname=test",
			expectedStatus: http.StatusUnauthorized,
			headers:        map[string]string{"Cookie": "token=wrongtoken"},
		},
		{
			name:           "Test downloading file rejects faulty request",
			method:         http.MethodGet,
			endpoint:       "/api/admin/get_file?machine=wrong",
			expectedStatus: http.StatusBadRequest,
			headers:        map[string]string{"Cookie": "token=example_token"},
		},
		{
			name:           "Test listing files is autheticated",
			method:         http.MethodGet,
			endpoint:       "/api/admin/list_files?machine=notauth",
			expectedStatus: http.StatusUnauthorized,
			headers:        map[string]string{"Cookie": "token=wrongtoken"},
		},
		{
			name:           "Test removing files is autheticated",
			method:         http.MethodPost,
			endpoint:       "/api/admin/remove_file",
			expectedStatus: http.StatusUnauthorized,
			body:           []byte(`{"hostname":"bogus", "filename":"bogus"}`),
			headers:        map[string]string{"Cookie": "token=wrongtoken"},
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

			assert.Equal(t, tc.expectedStatus, recorder.Code, tc.name)
		})
	}
}
