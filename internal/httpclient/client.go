package httpclient

import (
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: 500 * time.Millisecond,
}
