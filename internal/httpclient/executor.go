package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/IsmailCLN/tapir/internal/assert"
	"github.com/IsmailCLN/tapir/internal/parser"
	"github.com/IsmailCLN/tapir/internal/report"
)

// Optimize edilmiş shared HTTP client
var sharedClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
	Timeout: 10 * time.Second,
}

func RunAllTests(suite parser.TestSuite) {
	for _, tc := range suite.Tests {
		var body io.Reader
		if tc.Body != "" {
			body = bytes.NewBuffer([]byte(tc.Body))
		}

		req, err := http.NewRequest(tc.Method, tc.URL, body)
		if err != nil {
			report.PrintResult(tc.Name, assert.AssertionResult{
				Pass:        false,
				Description: "request creation error: " + err.Error(),
				Expected:    "valid request",
				Actual:      err.Error(),
			}, 0, 0, tc.Expect.Status, 0)
			continue
		}

		for k, v := range tc.Headers {
			req.Header.Set(k, v)
		}

		start := time.Now()
		resp, err := sharedClient.Do(req)
		duration := time.Since(start).Milliseconds()

		if err != nil {
			report.PrintResult(tc.Name, assert.AssertionResult{
				Pass:        false,
				Description: "HTTP error: " + err.Error(),
				Expected:    "success",
				Actual:      err.Error(),
			}, duration, 0, tc.Expect.Status, 0)
			continue
		}
		defer resp.Body.Close()

		// Body uzunluğunu hesapla ama içeriği işlemiyoruz (Discard)
		size, _ := io.Copy(io.Discard, resp.Body)

		result := assert.AssertStatus(tc.Expect.Status, resp.StatusCode)

		report.PrintResult(
			tc.Name,
			result,
			duration,
			size,
			tc.Expect.Status,
			resp.StatusCode,
		)
	}

	report.RenderResults()
}
