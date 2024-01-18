package femto

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Post data to the given URL
func Post[T any](url string, data T) error {
	// TODO:
	// - Logging
	// - Customizable transport

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	resp.Body.Close()

	return err
}
