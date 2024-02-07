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

type HTTPResponseContent[T any] struct {
	Headers map[string]string
	Body    T
	Status  int
}

type PostFunc[T any] func(T, *http.Request, *slog.Logger) (HTTPResponseContent[types.Unit], error)
type GetFunc[T any] func(*http.Request, *slog.Logger) (HTTPResponseContent[T], error)

func (femto *Femto) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	femto.mux.ServeHTTP(w, r)
}

func OnPost[T any](femto *Femto, pattern string, handle PostFunc[T]) {
	femto.mux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		doPost(femto, writer, request, handle)
	})
}

func OnGet[T any](femto *Femto, pattern string, handle GetFunc[T]) {
	femto.mux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		doGet(femto, writer, request, handle)
	})
}

func FailHandler[T any](e error) (HTTPResponseContent[T], error) {
	return HTTPResponseContent[T]{}, e
}

func generateErrorLogger(log *slog.Logger, writer http.ResponseWriter) func(ctx string, status int, e error) {
	return func(ctx string, status int, e error) {
		log.Error(ctx, "err", e)
		http.Error(writer, e.Error(), status)
	}
}

func doGet[T any](femto *Femto, writer http.ResponseWriter, request *http.Request, handle GetFunc[T]) {
	reqNo := femto.nextReqNo()
	log := femto.logger().With(slog.Uint64("req_no", reqNo))

	log.Info("New Request", "method", request.Method, "url", request.URL, "from", request.RemoteAddr)

	if request.Method != http.MethodGet {
		log.Info("Wanted GET, returned 405")
		http.Error(writer, "Expected GET", http.StatusMethodNotAllowed)
		return
	}

	data, err := handle(request, log)

	if data.Status == 0 {
		data.Status = http.StatusInternalServerError
	}

	ise := generateErrorLogger(log, writer)

	if err != nil {
		ise("There was an error in the handler", data.Status, err)
		return
	}

	jsonb, err := json.Marshal(data.Body)
	if err != nil {
		ise("There was an error in trying to parse the handler's response into JSON", data.Status, err)
		return
	}

	for key, value := range data.Headers {
		writer.Header().Add(key, value)
	}

	_, err = writer.Write(jsonb)

	if err != nil {
		ise("There was an error in trying to respond to the user", data.Status, err)
	}
}

func doPost[T any](f *Femto, writer http.ResponseWriter, reader *http.Request, handle PostFunc[T]) {
	reqNo := f.nextReqNo()
	log := f.logger().With(slog.Uint64("req_no", reqNo))

	log.Info("New Request", "method", reader.Method, "url", reader.URL, "from", reader.RemoteAddr)

	if reader.Method != http.MethodPost {
		log.Info("Wanted POST, returned 405")
		http.Error(writer, "Expected POST", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Ensure all field's are present. I'm so mad this is hard.
	var reqData T
	decoder := json.NewDecoder(reader.Body)

	if err := decoder.Decode(&reqData); err != nil {
		log.Info("Failed to unmarshal JSON", "error", err)
		http.Error(writer, "Failed to decode the provided JSON"+err.Error(), http.StatusBadRequest)
		return
	}

	ise := generateErrorLogger(log, writer)

	data, userErr := handle(reqData, reader, log)

	if data.Status == 0 {
		data.Status = http.StatusInternalServerError
	}

	if userErr != nil {

		ise("There was an error in trying to authenticate the user", data.Status, userErr)
		return
	}

	for key, value := range data.Headers {
		writer.Header().Add(key, value)
	}

	_, err := writer.Write(make([]byte, 0))

	if err != nil {
		ise("There was an error in trying to respond to the user", data.Status, err)
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
