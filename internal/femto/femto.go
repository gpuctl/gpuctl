// Package femto provides a tiny abstraction layer for writing API's that serve JSON over HTTP.
package femto

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync/atomic"
)

type Femto struct {
	mux   http.ServeMux
	reqNo atomic.Uint64
}

func (f *Femto) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mux.ServeHTTP(w, r)
}

func OnPost[T any, R any](f *Femto, pattern string, handle PostFuncPure[T, R]) {
	f.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		doPost(f, w, r, handle)
	})
}

func OnGet[T any](f *Femto, pattern string, handle GetFunc[T]) {
	f.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		doGet(f, w, r, handle)
	})
}

// correctMethod returns wether the request had the given method.
//
// If not, it will return a suitable reponce: either No Content for an OPTIONS request,
// or
func correctMethod(method string, req *http.Request, w http.ResponseWriter, log *slog.Logger) bool {
	switch req.Method {
	case method:
		return true
	case http.MethodOptions:

		// We don't set Access-Control-Allow-Origin, as that'd done globally
		// so it can be enabled for dev, but not prod.

		// https://developer.mozilla.org/en-US/docs/Glossary/Preflight_request
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/OPTIONS

		w.Header().Set("Allow", method)
		w.Header().Set("Access-Control-Allow-Methods", method)
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return false
	default:
		log.Info("Wrong method, returning 405", "expected", method, "got", req.Method)
		http.Error(w, "Expected "+method, http.StatusMethodNotAllowed)
		return false
	}

}

func doGet[T any](f *Femto, w http.ResponseWriter, r *http.Request, handle GetFunc[T]) {
	reqNo := f.nextReqNo()
	log := f.logger().With(slog.Uint64("req_no", reqNo))

	log.Info("New Request", "method", r.Method, "url", r.URL, "from", r.RemoteAddr)

	if !correctMethod(http.MethodGet, r, w, log) {
		return
	}

	ise := func(ctx string, e error) {
		log.Error(ctx, "err", e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}

	data, err := handle(r, log)
	if err != nil {
		// Handle authentication related errors
		if errors.Is(err, NotAuthenticatedError) {
			log.Info("Invalid attempt at authentication")
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}
		ise("application error", err)
		return
	}

	jsonb, err := json.Marshal(data)
	if err != nil {
		ise("marshaling to json failed", err)
		return
	}

	_, err = w.Write(jsonb)
	if err != nil {
		ise("writing failed", err)
	}
}

func doPost[T any, R any](f *Femto, w http.ResponseWriter, r *http.Request, handle PostFuncPure[T, R]) {
	reqNo := f.nextReqNo()
	log := f.logger().With(slog.Uint64("req_no", reqNo))

	log.Info("New Request", "method", r.Method, "url", r.URL, "from", r.RemoteAddr)

	if !correctMethod(http.MethodPost, r, w, log) {
		return
	}

	// TODO: Ensure all field's are present. I'm so mad this is hard.
	var reqData T
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		log.Info("Failed to unmarshal JSON", "error", err)
		http.Error(w, "failed to decode json: "+err.Error(), http.StatusBadRequest)
		return
	}

	ise := func(ctx string, e error) {
		log.Error(ctx, "err", e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}

	data, userErr := handle(reqData, r, log)

	if userErr != nil {
		// Handle authentication related errors
		if errors.Is(userErr, NotAuthenticatedError) {
			log.Info("Invalid attempt at authentication")
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}
		// TODO: Nicer error
		log.Info("Error", "err", userErr)
		http.Error(w, userErr.Error(), http.StatusInternalServerError)
		return
	}

	jsonb, err := json.Marshal(data)
	if err != nil {
		ise("marshaling to json failed", err)
		return
	}

	_, err = w.Write(jsonb)
	if err != nil {
		ise("writing failed", err)
	}
}

func (f *Femto) nextReqNo() uint64 {
	return f.reqNo.Add(1)
}

func (f *Femto) logger() *slog.Logger {
	// TODO: Associate a logger with Femto.
	return slog.Default()
}

func PurePost[T any](data T, r *http.Request, l *slog.Logger) (struct{}, error) {
	return struct{}{}, nil
}

// go functional hackery
func ParallelCompose[T any, R any](base PostFunc[T], pure PostFuncPure[T, R]) PostFuncPure[T, R] {
	return func(data T, r *http.Request, l *slog.Logger) (R, error) {
		err := base(data, r, l)

		if pure != nil {
			ret, e := pure(data, r, l)
			return ret, errors.Join(err, e)
		}

		var zero R
		return zero, err
	}
}

func WrapPostFunc[T any](f PostFunc[T]) PostFuncPure[T, struct{}] {
	return func(data T, r *http.Request, l *slog.Logger) (struct{}, error) {
		return struct{}{}, f(data, r, l)
	}
}

func WrapGetFunc(f GetFuncUnit) GetFunc[struct{}] {
	return func(r *http.Request, l *slog.Logger) (struct{}, error) {
		return struct{}{}, f(r, l)
	}
}

type PostFuncPure[T any, R any] func(T, *http.Request, *slog.Logger) (R, error)
type PostFunc[T any] func(T, *http.Request, *slog.Logger) error
type GetFuncUnit func(*http.Request, *slog.Logger) error
type GetFunc[T any] func(*http.Request, *slog.Logger) (T, error)
