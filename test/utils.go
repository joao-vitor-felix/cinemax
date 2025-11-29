package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func MakeRequest(method, path string, body any) *http.Request {
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}
