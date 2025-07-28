package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeadersForwarded(t *testing.T) {
	var gotAuth string
	mux := http.NewServeMux()
	mux.HandleFunc("/h", func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(200)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL+"/h", nil)
	req.Header.Set("Authorization", "Bearer abc123")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()

	if gotAuth != "Bearer abc123" {
		t.Fatalf("header not forwarded, got %q", gotAuth)
	}
}
