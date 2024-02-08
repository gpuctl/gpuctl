package authentication

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/types"
)

var (
	NotAuthenticatedError   = errors.New("User does not have a valid authentication token")
	InvalidCredientalsError = errors.New("Invalid credientals for creating an authentication token")
)

type AuthToken = string

type Authenticator[AuthCredientals any] interface {
	CreateToken(AuthCredientals) (AuthToken, error)
	RevokeToken(AuthToken) error
	CheckToken(AuthToken) bool
}

func AuthWrapGet[A any, T any](auth Authenticator[A], handle femto.GetFunc[T]) femto.GetFunc[T] {
	return func(request *http.Request, logger *slog.Logger) (*femto.Response[T], error) {
		authorisation := request.Header.Get("Authorization")

		token, bearerPresent := strings.CutPrefix(authorisation, "Bearer ")

		if !bearerPresent {
			return &femto.Response[T]{Status: http.StatusUnauthorized}, NotAuthenticatedError
		}

		if !auth.CheckToken(token) {
			return &femto.Response[T]{Status: http.StatusUnauthorized}, NotAuthenticatedError
		}

		return handle(request, logger)
	}
}

func AuthWrapPost[A any, T any](auth Authenticator[A], handle femto.PostFunc[T]) femto.PostFunc[T] {
	return func(data T, request *http.Request, logger *slog.Logger) (*femto.Response[types.Unit], error) {
		authorisation := request.Header.Get("Authorization")

		token, bearerPresent := strings.CutPrefix(authorisation, "Bearer ")

		if !bearerPresent {
			return &femto.Response[types.Unit]{Status: http.StatusUnauthorized}, NotAuthenticatedError
		}

		if !auth.CheckToken(token) {
			return &femto.Response[types.Unit]{Status: http.StatusUnauthorized}, NotAuthenticatedError
		}
		return handle(data, request, logger)
	}
}
