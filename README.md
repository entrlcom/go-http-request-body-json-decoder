# HTTP request body JSON decoder

## Table of Content

- [Examples](#examples)
- [License](#license)

## Examples

```go
package main

import (
	"net/http"
	"time"

	"entrlcom.dev/go-http-request-body-json-decoder"
)

const maxBytes = 1 << (10 * 1) * 2 // 2 KiB.

type Request struct {
	DateOfBirth time.Time `json:"date_of_birth"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var request Request

	// Decode HTTP request body to struct.
	if err := http_request_body_json_decoder.Decode(w, r, &request, maxBytes); err != nil {
		// TODO: Handle error.
		return
	}

	// ...
}

```

## License

[MIT](https://choosealicense.com/licenses/mit/)
