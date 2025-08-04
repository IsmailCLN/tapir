package runner

import (
	"maps"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/IsmailCLN/tapir/internal/assert"
	"github.com/IsmailCLN/tapir/internal/domain"
	"github.com/IsmailCLN/tapir/internal/httpclient"
)

// Result holds the outcome of a single request-level expectation.
type Result struct {
	Suite    string // test-suite name
	Request  string // request name
	Passed   bool   // true when expectation succeeded
	Err      error  // populated when Passed == false
	TestName string
}

func Run(ctx context.Context, suites []domain.TestSuite) ([]Result, error) {
	var results []Result

	for _, s := range suites {
		for _, r := range s.Requests {

			// 1) HTTP isteğini inşa et
			req, err := http.NewRequest(r.Req.Method, r.Req.URL, nil)
			if err != nil {
				// → İstek baştan çöktü: her expect için ayrı sonuç üret
				appendRequestErrorResults(&results, s.Name, r, err)
				continue
			}

			// 2) Gönder
			resp, err := httpclient.Do(ctx, req)
			if err != nil {
				// → Ağ/bağlantı hatası: yine tüm expect’leri hatalı say
				appendRequestErrorResults(&results, s.Name, r, err)
				continue
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			// 3) Expectation’ları değerlendir
			for _, exp := range r.Expect {
				// ➊ Kullanıcı YAML’inden gelen parametreleri kopyala
				kwargs := make(map[string]any, len(exp.Kwargs)+1)
				maps.Copy(kwargs, exp.Kwargs)
				// ➋ Otomatik ek parametreler
				kwargs["status_code"] = resp.StatusCode

				f, ok := assert.Get(exp.Type)
				if !ok {
					results = append(results, Result{
						Suite:    s.Name,
						Request:  r.Name,
						Passed:   false,
						Err:      fmt.Errorf("unknown expectation %s", exp.Type),
						TestName: exp.Type,
					})
					continue
				}

				err := f(body, kwargs)
				results = append(results, Result{
					Suite:    s.Name,
					Request:  r.Name,
					Passed:   err == nil,
					Err:      err,
					TestName: exp.Type,
				})
			}
		}
	}
	return results, nil
}

func appendRequestErrorResults(res *[]Result, suite string, r domain.TestRequest, err error) {
	if len(r.Expect) == 0 {
		*res = append(*res, Result{
			Suite:    suite,
			Request:  r.Name,
			Passed:   false,
			Err:      err,
			TestName: "request_error",
		})
		return
	}
	for _, exp := range r.Expect {
		*res = append(*res, Result{
			Suite:    suite,
			Request:  r.Name,
			Passed:   false,
			Err:      err,
			TestName: exp.Type,
		})
	}
}
