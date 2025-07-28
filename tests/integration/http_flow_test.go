package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IsmailCLN/tapir/internal/assert"
)

func TestHTTPFlow_StatusAndBody(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "name": "Leanne"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/users/1")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()

	st := assert.AssertStatus(200, resp.StatusCode)
	if !st.Pass {
		t.Fatalf("status mismatch: %+v", st)
	}

	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	b, _ := json.Marshal(got) // normalize

	expected := `{"id":1,"name":"Leanne"}`
	bd := assert.AssertBody(expected, string(b))
	if !bd.Pass {
		t.Fatalf("body mismatch: %+v", bd)
	}
}
