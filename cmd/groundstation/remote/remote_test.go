package remote

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gpuctl/gpuctl/internal/status"
	"github.com/stretchr/testify/assert"
)

/* Status Object Construction */

func TestBuildStatusObject(t *testing.T) {
	validJSON := []byte(`{"gpu_name": "Test GPU", "gpu_brand": "BrandX", "driver_ver": "1.0", "memory_total": 4096, "memory_util": 50, "gpu_util": 50, "memory_used": 2048, "fan_speed": 70, "gpu_temp": 60}`)
	expectedPacket := status.GPUStatusPacket{
		Name:              "Test GPU",
		Brand:             "BrandX",
		DriverVersion:     "1.0",
		MemoryTotal:       4096,
		MemoryUtilisation: 50,
		GPUUtilisation:    50,
		MemoryUsed:        2048,
		FanSpeed:          70,
		Temp:              60,
	}

	packet, err := buildStatusObject(validJSON)
	assert.NoError(t, err)
	assert.Equal(t, expectedPacket, packet)

}

func TestBuildMalformedObject(t *testing.T) {
	invalidJSON := []byte(`{"gpu_name": "Test GPU", "gpu_brand": "BrandX", "driver_ver": 1.0}`)
	_, err := buildStatusObject(invalidJSON)

	assert.Error(t, err)
}

/* Submission Side-Effect Testing */

type corruptedReader struct{}

func (r corruptedReader) Read(p []byte) (int, error) {
	return 0, errors.New("Read corrupted")
}

func simulateCorruptedSubmissionAndGetResponse(method string) *http.Response {
	req := httptest.NewRequest(method, "/api/submit", &corruptedReader{})
	w := httptest.NewRecorder()

	HandleStatusSubmission(w, req)

	res := w.Result()
	defer res.Body.Close()

	return res
}

func simulateSubmissionAndGetResponse(submission []byte, method string) *http.Response {
	req := httptest.NewRequest(method, "/api/submit", bytes.NewBuffer(submission))
	w := httptest.NewRecorder()

	HandleStatusSubmission(w, req)

	res := w.Result()
	defer res.Body.Close()

	return res
}

func assertResponseHas(t *testing.T, res *http.Response, expectedStatusCode int, expectedMessage string) {
	assert.Equal(t, expectedStatusCode, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, expectedMessage, string(data))
}

func TestHandleStatusSubmission(t *testing.T) {
	validJSON := []byte(`{"gpu_name": "Test GPU", "gpu_brand": "BrandX", "driver_ver": "1.0", "memory_total": 4096, "memory_util": 50, "gpu_util": 50, "memory_used": 2048, "fan_speed": 70, "gpu_temp": 60}`)

	res := simulateSubmissionAndGetResponse(validJSON, "POST")

	assertResponseHas(t, res, http.StatusOK, "OK: Submission processed successfully")
}

func TestHandleWrongMethodSubmission(t *testing.T) {
	validJSON := []byte(`{"gpu_name": "Test GPU", "gpu_brand": "BrandX", "driver_ver": "1.0", "memory_total": 4096, "memory_util": 50, "gpu_util": 50, "memory_used": 2048, "fan_speed": 70, "gpu_temp": 60}`)

	res := simulateSubmissionAndGetResponse(validJSON, "GET")

	assertResponseHas(t, res, http.StatusBadRequest, "Invalid method for status submission\n")
}

func TestHandleBadJsonSubmission(t *testing.T) {
	res := simulateCorruptedSubmissionAndGetResponse("POST")

	assertResponseHas(t, res, http.StatusBadRequest, "Malformed request body detected\n")
}

func TestHandleBadJsonDeserialisation(t *testing.T) {
	invalidJSON := []byte(`{"gpu_name": "Test GPU", "gpu_brand": "BrandX", "driver_ver": 1.0}`)

	res := simulateSubmissionAndGetResponse(invalidJSON, "POST")

	assertResponseHas(t, res, http.StatusBadRequest, "JSON deserialisation was not successful\n")
}

// Once DB connection is made, should test that we can fail on error states during
// interaction with DB
