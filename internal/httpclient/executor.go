package httpclient

import (
	"bytes"
	"io"
	"net/http"

	"github.com/IsmailCLN/tapir/internal/assert"
	"github.com/IsmailCLN/tapir/internal/parser"
	"github.com/IsmailCLN/tapir/internal/report"
)

func RunAllTests(suite parser.TestSuite) {
	for _, tc := range suite.Tests {
		client := &http.Client{}
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
			})
			continue
		}

		for k, v := range tc.Headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			report.PrintResult(tc.Name, assert.AssertionResult{
				Pass:        false,
				Description: "HTTP error: " + err.Error(),
				Expected:    "success",
				Actual:      err.Error(),
			})
			continue
		}
		defer resp.Body.Close()

		result := assert.AssertStatus(tc.Expect.Status, resp.StatusCode)
		report.PrintResult(tc.Name, result)
	}

	report.RenderResults()
}
