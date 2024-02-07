package femto

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

var (
	NotAuthenticatedError   = errors.New("User does not have a valid authentication token")
	InvalidCredientalsError = errors.New("Invalid credientals for creating an authentication token")
)

const (
	AuthCookieField = "auth_cookie"
)

type AuthToken = string

type Authenticator[AuthCredientals any] interface {
	CreateToken(AuthCredientals) (AuthToken, error)
	RevokeToken(AuthToken) error
	CheckToken(AuthToken) bool
}

func AuthWrapGet[A any, T any](auth Authenticator[A], handle GetFunc[T]) GetFunc[T] {
	return func(rr *http.Request, ll *slog.Logger) (T, error) {
		f := AuthWrapPost(auth, func(zero struct{}, r *http.Request, l *slog.Logger) (T, error) {
			return handle(r, l)
		})
		return f(struct{}{}, rr, ll)
	}
}

func AuthWrapPost[A any, T any, R any](auth Authenticator[A], handle PostFuncPure[T, R]) PostFuncPure[T, R] {
	return func(data T, r *http.Request, l *slog.Logger) (R, error) {
		var zero R
		authorisation := r.Header.Get("Authorization")
		split := strings.Split(authorisation, "Bearer ")

		res := false
		if len(split) > 1 {
			res = auth.CheckToken(split[1])
		}

		if !res {
			return zero, NotAuthenticatedError
		}
		return handle(data, r, l)
	}
}
