package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const maxBodySize = 1 << 20

var (
	ErrEmptyBody     = errors.New("request body is empty")
	ErrTooLarge      = errors.New("request body is too large")
	ErrInvalidJSON   = errors.New("invalid JSON format")
	ErrUnknownFields = errors.New("request contains unknown fields")
)

func ToJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	var b []byte
	var err error

	if os.Getenv("DEBUG") == "true" {
		b, err = json.MarshalIndent(data, "", "  ")
	} else {
		b, err = json.Marshal(data)
	}

	if err != nil {
		http.Error(w, `{"error":"failed to encode json"}`, http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func FromJSON(r *http.Request, dest any) error {
	if dest == nil {
		return errors.New("dest is nil")
	}

	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		if err.Error() == "http: request body too large" {
			return ErrTooLarge
		}
		return err
	}

	if len(body) == 0 {
		return ErrEmptyBody
	}

	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.DisallowUnknownFields()

	if err := dec.Decode(dest); err != nil {
		if strings.Contains(err.Error(), "unknown field") {
			return ErrUnknownFields
		}
		return ErrInvalidJSON
	}

	sanitize(dest)

	return nil
}

func sanitize(v any) {
	val := reflect.ValueOf(v).Elem()

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)

		if f.Kind() == reflect.String && f.CanSet() {
			f.SetString(strings.TrimSpace(f.String()))
		}
	}
}
