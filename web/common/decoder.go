package common

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dst); err != nil {
		return ErrInvalidJSON
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return ErrInvalidJSON
	}
	return nil
}
