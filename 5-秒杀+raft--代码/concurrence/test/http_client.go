package test

import (
	"net/http"
	"testing"
)

func TestHttpClient(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		http.Get("http://127.0.0.1:5678")
	}
}
