package femto

import (
	"errors"
	"log/slog"
	"net/http"
)

var (
	NotAuthenticatedError = errors.New("User does not have a valid authentication token")
)

type AuthToken = string

type Authenticator interface {
	Authenticate(username string, password string) (AuthToken, error)
	Unauthenticate(token AuthToken) error
	CheckAuth(token AuthToken) bool
}

func AuthWrapGet[T any](auth Authenticator, handle GetFunc[T]) GetFunc[T] {
	return func(r *http.Request, l *slog.Logger) (T, error) {
		var zero T
		cookie, err := r.Cookie("auth_cookie")
		if err != nil {
			return zero, err
		}

		res := auth.CheckAuth(cookie.String())
		if !res {
			return zero, NotAuthenticatedError
		}
		return handle(r, l)
	}
}

func AuthWrapPost[T any](auth Authenticator, handle PostFunc[T]) PostFunc[T] {
	return func(data T, r *http.Request, l *slog.Logger) error {
		cookie, err := r.Cookie("auth_cookie")
		if err != nil {
			return err
		}

		res := auth.CheckAuth(cookie.String())
		if !res {
			return NotAuthenticatedError
		}
		return handle(data, r, l)
	}
}
