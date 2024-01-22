package femto_test

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/stretchr/testify/assert"
)

var _ (http.Handler) = (*femto.Femto)(nil)

func TestError404(t *testing.T) {
	t.Parallel()

	mux := new(femto.Femto)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestWrongMethod(t *testing.T) {
	t.Parallel()

	mux := new(femto.Femto)

	femto.OnPost[struct{}](mux, "/postme", func(s struct{}, r *http.Request, l *slog.Logger) error {
		return nil
	})

	req := httptest.NewRequest("GET", "/postme", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
}

func TestNotJson(t *testing.T) {
	t.Parallel()

	mux := new(femto.Femto)

	femto.OnPost[struct{}](mux, "/postme", func(s struct{}, r *http.Request, l *slog.Logger) error {
		return nil
	})

	req := httptest.NewRequest("POST", "/postme", bytes.NewBufferString("not json at all"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, w.Code, http.StatusBadRequest)

	data, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "failed to decode json")
}

func TestUserError(t *testing.T) {
	t.Parallel()

	mux := new(femto.Femto)
	femto.OnPost[struct{}](mux, "/postme", func(s struct{}, r *http.Request, l *slog.Logger) error {
		return errors.New("their is no spoon")
	})

	req := httptest.NewRequest("POST", "/postme", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, w.Code, http.StatusInternalServerError)

	data, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "their is no spoon")
}
