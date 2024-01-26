package groundstation_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gpuctl/gpuctl/internal/groundstation"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"github.com/stretchr/testify/assert"
)

/* Submission Side-Effect Testing */

type corruptedReader struct{}

func (r corruptedReader) Read(p []byte) (int, error) {
	return 0, errors.New("Read corrupted")
}

func simulateCorruptedSubmissionAndGetResponse(method string) *http.Response {
	req := httptest.NewRequest(method, uplink.GPUStatsUrl, &corruptedReader{})
	w := httptest.NewRecorder()

	s := groundstation.NewServer()
	s.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	return res
}

func simulateSubmissionAndGetResponse(submission []byte, method string) *http.Response {
	req := httptest.NewRequest(method, uplink.GPUStatsUrl, bytes.NewBuffer(submission))
	w := httptest.NewRecorder()

	s := groundstation.NewServer()
	s.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	return res
}

func assertResponseHas(t *testing.T, res *http.Response, expectedStatusCode int, expectedMessage string) {
	t.Helper()

	assert.Equal(t, expectedStatusCode, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Contains(t, string(data), expectedMessage)
}

func TestHandleStatusSubmission(t *testing.T) {
	t.Parallel()
	validJSON := []byte(`{"gpu_name": "Test GPU", "gpu_brand": "BrandX", "driver_ver": "1.0", "memory_total": 4096, "memory_util": 50, "gpu_util": 50, "memory_used": 2048, "fan_speed": 70, "gpu_temp": 60}`)

	res := simulateSubmissionAndGetResponse(validJSON, "POST")

	assertResponseHas(t, res, http.StatusOK, "OK")
}

func TestHandleWrongMethodSubmission(t *testing.T) {
	t.Parallel()
	validJSON := []byte(`{"gpu_name": "Test GPU", "gpu_brand": "BrandX", "driver_ver": "1.0", "memory_total": 4096, "memory_util": 50, "gpu_util": 50, "memory_used": 2048, "fan_speed": 70, "gpu_temp": 60}`)

	res := simulateSubmissionAndGetResponse(validJSON, "GET")

	assertResponseHas(t, res, http.StatusMethodNotAllowed, "Expected POST\n")
}

func TestHandleBadJsonSubmission(t *testing.T) {
	t.Parallel()
	res := simulateCorruptedSubmissionAndGetResponse("POST")

	assertResponseHas(t, res, http.StatusBadRequest, "Read corrupted\n")
}

func TestHandleBadJsonDeserialisation(t *testing.T) {
	t.Parallel()
	invalidJSON := []byte(`{"gpu_name": "Test GPU", "gpu_brand": "BrandX", "driver_ver": 1.0}`)

	res := simulateSubmissionAndGetResponse(invalidJSON, "POST")

	assertResponseHas(t, res, http.StatusBadRequest, "failed to decode json")
}

// Once DB connection is made, should test that we can fail on error states during
// interaction with DB
