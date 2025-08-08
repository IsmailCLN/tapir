package runner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"maps"

	"github.com/IsmailCLN/tapir/internal/assert"
	"github.com/IsmailCLN/tapir/internal/domain"
	"github.com/IsmailCLN/tapir/internal/httpclient"
	"github.com/IsmailCLN/tapir/internal/sharedcontext"
)

// Result holds the outcome of a single request-level expectation.
type Result struct {
	Suite    string
	Request  string
	Passed   bool
	Err      error
	TestName string
}

func Run(ctx context.Context, suites []domain.TestSuite) ([]Result, error) {
	var results []Result

	shared := sharedcontext.New()
	assert.SetSharedContext(shared)

	for _, s := range suites {
		for _, r := range s.Requests {

			// ----- 1. Build request body (string only for now) -----
			var bodyReader io.Reader
			if bodyStr, ok := r.Req.Body.(string); ok && bodyStr != "" {
				bodyReader = strings.NewReader(bodyStr)
			}

			// ----- 2. Construct HTTP request -----
			req, err := http.NewRequest(r.Req.Method, r.Req.URL, bodyReader)
			if err != nil {
				appendRequestErrorResults(&results, s.Name, r, err)
				continue
			}

			// ----- 3. Apply headers with placeholder substitution -----
			for k, v := range r.Req.Headers {
				if strings.Contains(v, "${token}") {
					if t, ok := shared.Get("token"); ok {
						v = strings.ReplaceAll(v, "${token}", t)
					}
				}
				req.Header.Set(k, v)
			}

			// ----- 4. Send request -----
			resp, err := httpclient.Do(ctx, req)
			if err != nil {
				appendRequestErrorResults(&results, s.Name, r, err)
				continue
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				appendRequestErrorResults(&results, s.Name, r, err)
				continue
			}

			// ----- 5. Evaluate expectations -----
			for _, exp := range r.Expect {
				// 5a. Copy user‑provided kwargs
				kwargs := make(map[string]any, len(exp.Kwargs)+1)
				maps.Copy(kwargs, exp.Kwargs)

				// 5b. Inject auto params
				kwargs["status_code"] = resp.StatusCode
				kwargs["headers"] = resp.Header

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

				err := f(bodyBytes, kwargs)
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
