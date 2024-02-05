// Package femto provides a tiny abstraction layer for writing API's that serve JSON over HTTP.
package femto

import (
	"encoding/json"
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

func OnPost[T any, R any](f *Femto, pattern string, handle PostFunc[T, R]) {
	f.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		doPost(f, w, r, handle)
	})
}

func OnGet[T any](f *Femto, pattern string, handle GetFunc[T]) {
	f.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		doGet(f, w, r, handle)
	})
}

func doGet[T any](f *Femto, w http.ResponseWriter, r *http.Request, handle GetFunc[T]) {
	reqNo := f.nextReqNo()
	log := f.logger().With(slog.Uint64("req_no", reqNo))

	log.Info("New Request", "method", r.Method, "url", r.URL, "from", r.RemoteAddr)

	if r.Method != http.MethodGet {
		log.Info("Wanted GET, returned 405")
		http.Error(w, "Expected GET", http.StatusMethodNotAllowed)
		return
	}

	ise := func(ctx string, e error) {
		log.Error(ctx, "err", e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}

	data, err := handle(r, log)
	if err != nil {
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

func doPost[T any, R any](f *Femto, w http.ResponseWriter, r *http.Request, handle PostFunc[T, R]) {
	reqNo := f.nextReqNo()
	log := f.logger().With(slog.Uint64("req_no", reqNo))

	log.Info("New Request", "method", r.Method, "url", r.URL, "from", r.RemoteAddr)

	if r.Method != http.MethodPost {
		log.Info("Wanted POST, returned 405")
		http.Error(w, "Expected POST", http.StatusMethodNotAllowed)
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

type PostFunc[T any, R any] func(T, *http.Request, *slog.Logger) (R, error)
type GetFunc[T any] func(*http.Request, *slog.Logger) (T, error)
