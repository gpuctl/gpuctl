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

func (f *Femto) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, f)
}

func (f *Femto) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mux.ServeHTTP(w, r)
}

func OnPost[T any](f *Femto, pattern string, handle HandlerFunc[T]) {
	f.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		doPost(f, w, r, handle)
	})
}

func doPost[T any](f *Femto, w http.ResponseWriter, r *http.Request, handle HandlerFunc[T]) {

	reqNo := f.nextReqNo()
	log := f.logger().With(slog.Uint64("req_no", reqNo))

	log.Info("New Request", "method", r.Method, "url", r.URL, "from", r.RemoteAddr)

	if r.Method != http.MethodPost {
		log.Info("Wanted GET, returned 405")
		http.Error(w, "Expected GET", http.StatusMethodNotAllowed)
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

	userErr := handle(reqData, r, log)
	if userErr != nil {
		// TODO: Nicer error
		log.Info("Error", "err", userErr)
		http.Error(w, userErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK"))
}

func (f *Femto) nextReqNo() uint64 {
	return f.reqNo.Add(1)
}

func (f *Femto) logger() *slog.Logger {
	// TODO: Associate a logger with Femto.
	return slog.Default()
}

type HandlerFunc[T any] func(T, *http.Request, *slog.Logger) error
