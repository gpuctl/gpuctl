// Package femto provides a tiny abstraction layer for writing API's that serve JSON over HTTP.
package femto

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync/atomic"

	"github.com/gpuctl/gpuctl/internal/types"
)

type Femto struct {
	mux   http.ServeMux
	reqNo atomic.Uint64
}

type Response[T any] struct {
	Headers map[string]string
	Body    T
	Status  int
}

type EmptyBodyResponse = Response[types.Unit]

type PostFunc[T any] func(T, *http.Request, *slog.Logger) (*EmptyBodyResponse, error)
type GetFunc[T any] func(*http.Request, *slog.Logger) (*Response[T], error)

func (femto *Femto) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	femto.mux.ServeHTTP(w, r)
}

func OnPost[T any](f *Femto, pattern string, handle PostFunc[T]) {
	f.mux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		doPost(f, writer, request, handle)
	})
}

func OnGet[T any](f *Femto, pattern string, handle GetFunc[T]) {
	f.mux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		doGet(f, writer, request, handle)
	})
}

func Ok[T any](content T) Response[T] {
	return Response[T]{Status: http.StatusAccepted, Body: content}
}

func generateErrorLogger(l *slog.Logger, w http.ResponseWriter) func(ctx string, status int, e error) {
	return func(ctx string, status int, e error) {
		l.Error(ctx, "err", e)
		http.Error(w, e.Error(), status)
	}
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

	data, err := handle(r, log)

	if data == nil {
		data = &Response[T]{}
	}

	if data.Status == 0 {
		data.Status = http.StatusInternalServerError
	}

	ise := generateErrorLogger(log, w)

	if err != nil {
		ise("There was an error in the server when trying to handle the provided request", data.Status, err)
		return
	}

	jsonb, err := json.Marshal(data.Body)
	if err != nil {
		ise("There was an error in trying to serialise the handler's response into JSON", data.Status, err)
		return
	}

	for key, value := range data.Headers {
		w.Header().Add(key, value)
	}

	_, err = w.Write(jsonb)

	if err != nil {
		ise("There was an error in trying to write to the user", data.Status, err)
	}
}

func doPost[T any](f *Femto, w http.ResponseWriter, r *http.Request, handle PostFunc[T]) {
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
		http.Error(w, "Failed to decode the provided JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	ise := generateErrorLogger(log, w)

	data, userErr := handle(reqData, r, log)

	if data == nil {
		data = &EmptyBodyResponse{}
	}

	if data.Status == 0 {
		data.Status = http.StatusInternalServerError
	}

	if userErr != nil {

		ise("There was an error in the server when trying to handle the provided request", data.Status, userErr)
		return
	}

	for key, value := range data.Headers {
		w.Header().Add(key, value)
	}

	_, err := w.Write(make([]byte, 0))

	if err != nil {
		ise("There was an error in trying to write to the user", data.Status, err)
		return
	}
}

func (f *Femto) nextReqNo() uint64 {
	return f.reqNo.Add(1)
}

func (f *Femto) logger() *slog.Logger {
	// TODO: Associate a logger with Femto.
	return slog.Default()
}
