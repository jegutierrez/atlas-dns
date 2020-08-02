package app_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jegutierrez/atlas-dns/internal/app"
)

func TestAppRouting(t *testing.T) {
	var (
		validDatabankLocationRequest = `{
		"x": "123.12",
		"y": "456.56",
		"z": "789.89",
		"vel": "20.0"
		}`
		invalidDatabankLocationRequest = `{
		"a": "123.12",
		"b": "456.56",
		"c": "789.89",
		"vel": "20.0"
		}`
		validPingResponse = "pong"
	)
	tt := []struct {
		description  string
		path         string
		method       string
		body         *string
		statusCode   int
		responseBody *string
	}{
		{"test ping", "/ping", http.MethodGet, nil, http.StatusOK, &validPingResponse},
		{"test valid databank location service", "/calculate-databank-location", http.MethodPost, &validDatabankLocationRequest, http.StatusOK, nil},
		{"test valid databank location service", "/calculate-databank-location", http.MethodPost, &invalidDatabankLocationRequest, http.StatusBadRequest, nil},
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			app, err := app.NewApp()
			defer app.Shutdown()
			if err != nil {
				t.Fatal("could not create app")
			}
			router := app.RouterSetup()
			var body io.Reader
			if tc.body != nil {
				body = bytes.NewReader([]byte(*tc.body))
			}
			req, err := http.NewRequest(tc.method, tc.path, body)
			if err != nil {
				t.Fatalf("could not send request to %s: %v", tc.path, err)
			}
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tc.statusCode {
				t.Errorf("unexpected status code for %s, want: %d, got: %d", tc.path, resp.Code, tc.statusCode)
			}
			if tc.responseBody != nil && *tc.responseBody != resp.Body.String() {
				t.Errorf("unexpected body for %s, want: %s, got: %s", tc.path, resp.Body.String(), *tc.responseBody)
			}
		})
	}
}
