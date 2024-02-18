package authentication

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/femto"
)

const TokenCookieName = "token"

var (
	NotAuthenticatedError   = errors.New("User does not have a valid authentication token")
	InvalidCredientalsError = errors.New("Invalid credientals for creating an authentication token")
)

// These would probably be safer as newtypes, but exactly where to use the raw
// types and where to use the newtypes is slightly subtle (i.e. should packets
// from front-end use raw types because those values haven't been checked yet?
// What about packets we send back to the front-end?) so I would like someone
// else to worry about refactoring this - NB
type AuthToken = string
type Username = string

type Authenticator[AuthCredientals any] interface {
	CreateToken(AuthCredientals) (AuthToken, error)
	RevokeToken(AuthToken) error
	// Returns the username associated with the authentication token if it is
	// valid, otherwise an error
	CheckToken(AuthToken) (Username, error)
}

func AuthWrapGet[A any, T any](auth Authenticator[A], handle femto.GetFunc[T]) femto.GetFunc[T] {
	return func(request *http.Request, logger *slog.Logger) (*femto.Response[T], error) {
		c, err := request.Cookie(TokenCookieName)
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
		c, err := request.Cookie(TokenCookieName)
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
