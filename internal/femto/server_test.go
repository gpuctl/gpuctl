package femto_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gpuctl/gpuctl/internal/authentication"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/types"
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

	femto.OnPost(mux, "/postme", func(s struct{}, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return &femto.EmptyBodyResponse{}, nil
	})

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

	femto.OnPost(mux, "/postme", func(s struct{}, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return &femto.EmptyBodyResponse{}, nil
	})

	req := httptest.NewRequest("POST", "/postme", bytes.NewBufferString("not json at all"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, w.Code, http.StatusBadRequest)

	data, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "Failed to decode the provided JSON")
}

func TestUserError(t *testing.T) {
	t.Parallel()

	mux := new(femto.Femto)
	femto.OnPost(mux, "/postme", func(s struct{}, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return nil, errors.New("their is no spoon")
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

	femto.OnGet(mux, "/happy", func(r *http.Request, l *slog.Logger) (*femto.Response[Foo], error) {
		return &femto.Response[Foo]{Body: Foo{101}}, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/happy", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, http.StatusOK, w.Code)

	j := w.Body.Bytes()
	assert.JSONEq(t, `{"X": 101}`, string(j))
}

func TestGetWrongMethod(t *testing.T) {
	t.Parallel()

	mux := new(femto.Femto)
	femto.OnGet(mux, "/api", func(r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return &femto.EmptyBodyResponse{}, nil
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
	femto.OnGet(mux, "/api", func(r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return &femto.EmptyBodyResponse{}, fmt.Errorf("Some application error")
	})

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	defer req.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Some application error")
}

type TestAuthenticator struct{}

func (auth TestAuthenticator) CreateToken(unit types.Unit) (authentication.AuthToken, error) {
	return "token", nil
}

func (auth TestAuthenticator) RevokeToken(token authentication.AuthToken) error {
	return nil
}

func (auth TestAuthenticator) CheckToken(token authentication.AuthToken) bool {
	return token == "token"
}

func TestValidAuthentication(t *testing.T) {
	t.Parallel()
	mux := new(femto.Femto)

	auth := TestAuthenticator{}

	/* ------ Get Handler ------ */

	getHandler :=
		func(r *http.Request, l *slog.Logger) (*femto.Response[string], error) {
			return &femto.Response[string]{Body: "OKGET"}, nil
		}

	authenticatedGetHandler :=
		authentication.AuthWrapGet[types.Unit, string](auth, getHandler)

	// Set up authenticated endpoints
	femto.OnGet(mux, "/auth", authenticatedGetHandler)

	/* ------- Post Handler ------- */

	postHandler := func(incoming types.Unit, request *http.Request, logger *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return &femto.EmptyBodyResponse{}, nil
	}

	authenticatedPostHandler := authentication.AuthWrapPost[types.Unit, types.Unit](auth, postHandler)

	femto.OnPost(mux, "/auth-post", authenticatedPostHandler)

	/* ---------- Test --------- */

	w := httptest.NewRecorder()

	// Set up authorisation
	reqGet := httptest.NewRequest("GET", "/auth", strings.NewReader("{}"))
	defer reqGet.Body.Close()
	reqPost := httptest.NewRequest("POST", "/auth-post", strings.NewReader("{}"))
	defer reqPost.Body.Close()
	reqGet.Header.Set("Authorization", "Bearer token")
	reqPost.Header.Set("Authorization", "Bearer token")

	mux.ServeHTTP(w, reqGet)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "OKGET")

	mux.ServeHTTP(w, reqPost)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInvalidAuthentication(t *testing.T) {
	t.Parallel()
	mux := new(femto.Femto)
	auth := TestAuthenticator{}

	/* ------ Get Handler ------ */

	getHandler :=
		func(r *http.Request, l *slog.Logger) (*femto.Response[string], error) {
			return &femto.Response[string]{Body: "OKGET"}, nil
		}

	authenticatedGetHandler :=
		authentication.AuthWrapGet[types.Unit, string](auth, getHandler)

	// Set up authenticated endpoints
	femto.OnGet(mux, "/auth", authenticatedGetHandler)

	/* ------- Post Handler ------- */

	postHandler := func(incoming types.Unit, request *http.Request, logger *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return &femto.EmptyBodyResponse{}, nil
	}

	authenticatedPostHandler := authentication.AuthWrapPost[types.Unit, types.Unit](auth, postHandler)

	femto.OnPost(mux, "/auth-post", authenticatedPostHandler)

	w := httptest.NewRecorder()

	// Set up requests
	reqGet := httptest.NewRequest("GET", "/auth", strings.NewReader("{}"))
	defer reqGet.Body.Close()
	reqPost := httptest.NewRequest("POST", "/auth-post", strings.NewReader("{}"))
	defer reqPost.Body.Close()
	reqGet.Header.Set("Authorization", "Bearer wrong")
	reqPost.Header.Set("Authorization", "Bearer wrong")

	mux.ServeHTTP(w, reqGet)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.NotContains(t, w.Body.String(), "OKGET")

	mux.ServeHTTP(w, reqPost)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestOk(t *testing.T) {
	{
		resp, _ := femto.Ok("hello")
		if resp.Status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Status)
		}
		if resp.Body != "hello" {
			t.Errorf("Expected body %q, got %q", "hello", resp.Body)
		}
	}

	{
		resp, _ := femto.Ok(123)
		if resp.Status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Status)
		}
		if resp.Body != 123 {
			t.Errorf("Expected body %d, got %d", 123, resp.Body)
		}
	}

	type CustomStruct struct {
		Name string
		Age  int
	}
	{
		expectedBody := CustomStruct{Name: "John", Age: 30}
		resp, _ := femto.Ok(expectedBody)
		if resp.Status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Status)
		}
		if resp.Body != expectedBody {
			t.Errorf("Expected body %+v, got %+v", expectedBody, resp.Body)
		}
	}
}
