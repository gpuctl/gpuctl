package authentication

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/femto"
)

var (
	NotAuthenticatedError   = errors.New("User does not have a valid authentication token")
	InvalidCredientalsError = errors.New("Invalid credientals for creating an authentication token")
)

type AuthToken = string

type Authenticator[AuthCredientals any] interface {
	CreateToken(AuthCredientals) (AuthToken, error)
	RevokeToken(AuthToken) error
	CheckToken(AuthToken) (string, error)
}

func AuthWrapGet[A any, T any](auth Authenticator[A], handle femto.GetFunc[T]) femto.GetFunc[T] {
	return func(request *http.Request, logger *slog.Logger) (*femto.Response[T], error) {
		c, err := request.Cookie("token")
		if err != nil {
			return &femto.Response[T]{Status: http.StatusUnauthorized}, NotAuthenticatedError
		}

		token := c.Value
		_, err = auth.CheckToken(token)
		if err != nil {
			return &femto.Response[T]{Status: http.StatusUnauthorized}, NotAuthenticatedError
		}

		return handle(request, logger)
	}
}

func AuthWrapPost[A any, T any](auth Authenticator[A], handle femto.PostFunc[T]) femto.PostFunc[T] {
	return func(data T, request *http.Request, logger *slog.Logger) (*femto.EmptyBodyResponse, error) {
		c, err := request.Cookie("token")
		if err != nil {
			return &femto.EmptyBodyResponse{Status: http.StatusUnauthorized}, NotAuthenticatedError
		}

		token := c.Value
		_, err = auth.CheckToken(token)
		if err != nil {
			return &femto.EmptyBodyResponse{Status: http.StatusUnauthorized}, NotAuthenticatedError
		}
		return handle(data, request, logger)
	}
}
