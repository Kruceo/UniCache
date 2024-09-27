package memory

import "net/http"

type Cache struct {
	Data      []byte
	Headers   http.Header
	Timestamp int
	Access    int
}
