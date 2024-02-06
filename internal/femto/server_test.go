package femto_test

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
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

	femto.OnPost(mux, "/postme", femto.WrapPostFunc(func(s struct{}, r *http.Request, l *slog.Logger) error {
		return nil
	}))

	req := httptest.NewRequest("GET", "/postme", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	data, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "Expected POST")

	assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
}

func TestNotJson(t *testing.T) {
	t.Parallel()

	mux := new(femto.Femto)

	femto.OnPost(mux, "/postme", femto.WrapPostFunc(func(s struct{}, r *http.Request, l *slog.Logger) error {
		return nil
	}))

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
	femto.OnPost[struct{}, struct{}](mux, "/postme", func(s struct{}, r *http.Request, l *slog.Logger) (struct{}, error) {
		return struct{}{}, errors.New("their is no spoon")
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

func TestGetHappyPath(t *testing.T) {
	t.Parallel()

	mux := new(femto.Femto)

	type Foo struct {
		X int
	}

	femto.OnGet(mux, "/happy", func(r *http.Request, l *slog.Logger) (Foo, error) {
		return Foo{101}, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/happy", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, http.StatusOK, w.Code)

	j := w.Body.Bytes()
	assert.JSONEq(t, `{"X": 101}`, string(j))
}

type unit struct{}

func TestGetWrongMethod(t *testing.T) {
	t.Parallel()

	mux := new(femto.Femto)
	femto.OnGet(mux, "/api", func(r *http.Request, l *slog.Logger) (unit, error) {
		return unit{}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/api", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Contains(t, w.Body.String(), "Expected GET")
}

func TestGetApplicationErr(t *testing.T) {
	t.Parallel()

	mux := new(femto.Femto)
	femto.OnGet(mux, "/api", func(r *http.Request, l *slog.Logger) (unit, error) {
		return unit{}, errors.New("Some application Error")
	})

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Some application Error")
}

type TestAuthenticator struct{}

func (auth TestAuthenticator) CreateToken(unit struct{}) (femto.AuthToken, error) {
	return "token", nil
}

func (auth TestAuthenticator) RevokeToken(token femto.AuthToken) error {
	return nil
}

func (auth TestAuthenticator) CheckToken(token femto.AuthToken) bool {
	return token == "token"
}

func TestValidAuthentication(t *testing.T) {
	t.Parallel()
	mux := new(femto.Femto)

	auth := TestAuthenticator{}

	// Set up authenticated endpoints
	femto.OnGet(mux, "/auth", femto.AuthWrapGet(auth,
		func(r *http.Request, l *slog.Logger) (string, error) {
			return "OKGET", nil
		}))

	femto.OnPost(mux, "/auth-post", femto.AuthWrapPost(auth,
		func(s struct{}, r *http.Request, l *slog.Logger) (string, error) {
			return "OKPOST", nil
		}))

	w := httptest.NewRecorder()

	//http.SetCookie(w, &http.Cookie{Name: femto.AuthCookieField, Value: "token"})

	// Set up requests
	authCookie := http.Cookie{Name: femto.AuthCookieField, Value: "token"}

	reqGet := httptest.NewRequest("GET", "/auth", strings.NewReader("{}"))
	reqGet.AddCookie(&authCookie)
	defer reqGet.Body.Close()
	reqPost := httptest.NewRequest("POST", "/auth-post", strings.NewReader("{}"))
	reqPost.AddCookie(&authCookie)
	defer reqPost.Body.Close()

	mux.ServeHTTP(w, reqGet)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, string(w.Body.Bytes()), "OKGET")

	mux.ServeHTTP(w, reqPost)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, string(w.Body.Bytes()), "OKPOST")
}

func TestInvalidAuthentication(t *testing.T) {
	t.Parallel()
	mux := new(femto.Femto)
	auth := TestAuthenticator{}

	// Set up authenticated endpoints
	femto.OnGet(mux, "/auth", femto.AuthWrapGet(auth,
		func(r *http.Request, l *slog.Logger) (string, error) {
			return "OKGET", nil
		}))

	femto.OnPost(mux, "/auth-post", femto.AuthWrapPost(auth,
		func(s struct{}, r *http.Request, l *slog.Logger) (string, error) {
			return "OKPOST", nil
		}))

	w := httptest.NewRecorder()

	//http.SetCookie(w, &http.Cookie{Name: femto.AuthCookieField, Value: "token"})

	// Set up requests
	authCookie := http.Cookie{Name: femto.AuthCookieField, Value: "chicken"}

	reqGet := httptest.NewRequest("GET", "/auth", strings.NewReader("{}"))
	reqGet.AddCookie(&authCookie)
	defer reqGet.Body.Close()
	reqPost := httptest.NewRequest("POST", "/auth-post", strings.NewReader("{}"))
	reqPost.AddCookie(&authCookie)
	defer reqPost.Body.Close()

	mux.ServeHTTP(w, reqGet)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.NotContains(t, string(w.Body.Bytes()), "OKGET")

	mux.ServeHTTP(w, reqPost)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.NotContains(t, string(w.Body.Bytes()), "OKPOST")
}
