package httpe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// WriteJSON marshals struct and writes it to ResposneWriter.
func WriteJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// ReadJSON decodes JSON into v from reader.
// For invalid JSON returns an error with ErrMalformedJSON sentinel
func ReadJSON(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) || errors.Is(err, io.ErrUnexpectedEOF) {
			err = fmt.Errorf("%w: %v", ErrMalformedJSON, err)
		}
		return err
	}
	return nil
}
