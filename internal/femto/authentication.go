package femto

import (
	"errors"
	"log/slog"
	"net/http"
)

var (
	NotAuthenticatedError = errors.New("User does not have a valid authentication token")
	InvalidCredientalsError = errors.New("Invalid credientals for creating an authentication token")
)

const (
	AuthCookieField = "auth_cookie"
)

type AuthToken = string

type Authenticator[T any] interface {
	CreateToken(T) (AuthToken, error)
	RevokeToken(token AuthToken) error
	CheckToken(token AuthToken) bool
}

func AuthWrapGet[A any, T any](auth Authenticator[A], handle GetFunc[T]) GetFunc[T] {
	return func(r *http.Request, l *slog.Logger) (T, error) {
		var zero T
		cookie, err := r.Cookie(AuthCookieField)
		if err != nil {
			return zero, err
		}

		res := auth.CheckToken(cookie.Value)
		if !res {
			return zero, NotAuthenticatedError
		}
		return handle(r, l)
	}
}

func AuthWrapPost[A any, T any, R any](auth Authenticator[A], handle PostFunc[T, R]) PostFunc[T, R] {
	return func(data T, r *http.Request, l *slog.Logger) (R, error) {
		var zero R
		cookie, err := r.Cookie(AuthCookieField)
		if err != nil {
			return zero, err
		}

		res := auth.CheckToken(cookie.Value)
		if !res {
			return zero, NotAuthenticatedError
		}
		return handle(data, r, l)
	}
}
