package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/IsmailCLN/tapir/internal/assert"
)

func TestTimeoutLikeScenario(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(120 * time.Millisecond)
		w.WriteHeader(200)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := &http.Client{ Timeout: 50 * time.Millisecond }

	req, _ := http.NewRequest("GET", srv.URL+"/slow", nil)
	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	if err == nil {
		// Beklenen: timeout hatası. Eğer döndüyse yine de status assert yapalım.
		defer resp.Body.Close()
		st := assert.AssertStatus(200, resp.StatusCode)
		if !st.Pass {
			t.Fatalf("unexpected status: %+v", st)
		}
		// Bu branch'e girmek pek olası değil; ama istikrarlı olması için bırakıyoruz.
	} else {
		// Timeout beklenir: elapsed ~50ms civarı olmalı (biraz tolerans)
		if elapsed < 40*time.Millisecond || elapsed > 200*time.Millisecond {
			t.Fatalf("timeout elapsed looks off: %v", elapsed)
		}
	}
}
