package go_http_request_body_json_decoder

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

var (
	ErrInternal              = errors.New("internal error")
	ErrInvalidJSON           = errors.New("invalid json")
	ErrRequestEntityTooLarge = errors.New("request entity too large")
	ErrUnsupportedMediaType  = errors.New("unsupported media type")
)

func Decode(
	w http.ResponseWriter,
	r *http.Request,
	v any,
	n int64,
) error {
	contentType := r.Header.Get("Content-Type")
	if len(contentType) == 0 {
		return ErrUnsupportedMediaType
	}

	if strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0])) != "application/json" {
		return ErrUnsupportedMediaType
	}

	r.Body = http.MaxBytesReader(w, r.Body, n)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&v); err != nil {
		var (
			httpMaxBytesError      *http.MaxBytesError
			jsonSyntaxError        *json.SyntaxError
			jsonUnmarshalTypeError *json.UnmarshalTypeError
		)

		switch {
		case
			strings.HasPrefix(err.Error(), "json: unknown field "), // https://github.com/golang/go/issues/29035
			errors.Is(err, io.EOF),
			errors.Is(err, io.ErrUnexpectedEOF), // https://github.com/golang/go/issues/25956
			errors.As(err, &jsonSyntaxError),
			errors.As(err, &jsonUnmarshalTypeError):
			return errors.Join(err, ErrInvalidJSON)
		case errors.As(err, &httpMaxBytesError):
			return errors.Join(err, ErrRequestEntityTooLarge)
		default:
			return errors.Join(err, ErrInternal)
		}
	}

	if decoder.More() {
		return ErrInvalidJSON
	}

	return nil
}
