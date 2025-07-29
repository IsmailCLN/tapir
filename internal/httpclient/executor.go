package httpclient

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/IsmailCLN/tapir/internal/assert"
	"github.com/IsmailCLN/tapir/internal/config"
	"github.com/IsmailCLN/tapir/internal/parser"
	"github.com/IsmailCLN/tapir/internal/report"
)

var sharedTransport = &http.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 10,
	IdleConnTimeout:     90 * time.Second,
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	TLSHandshakeTimeout:   5 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var sharedClient = &http.Client{
	Transport: sharedTransport,
	Timeout:   10 * time.Second,
}

func executeSuite(suite parser.TestSuite) {
	// CLI’dan gelen timeout’u uygula
	sharedClient.Timeout = config.HTTPTimeout

abortedLoop:
	for _, tc := range suite.Tests {
		var body io.Reader
		if tc.Body != "" {
			body = bytes.NewBuffer([]byte(tc.Body))
		}

		req, err := http.NewRequest(tc.Method, tc.URL, body)
		if err != nil {
			report.PrintResult(tc.Name, "request", assert.AssertionResult{
				Pass:        false,
				Description: "request creation error: " + err.Error(),
				Expected:    "valid request",
				Actual:      err.Error(),
			}, 0, 0, 0, 0)

			if config.FailFast {
				break abortedLoop
			}
			continue
		}

		// YAML header’larını uygula
		for k, v := range tc.Headers {
			req.Header.Set(k, v)
		}

		// İsteği gönder
		ctx, cancel := context.WithTimeout(req.Context(), config.HTTPTimeout)
		req = req.WithContext(ctx)

		start := time.Now()
		resp, err := sharedClient.Do(req)
		duration := time.Since(start).Milliseconds()
		cancel()

		if err != nil {
			report.PrintResult(tc.Name, "http", assert.AssertionResult{
				Pass:        false,
				Description: "HTTP error: " + err.Error(),
				Expected:    "success",
				Actual:      err.Error(),
			}, duration, 0, tc.Expect.Status, 0)

			if config.FailFast {
				break abortedLoop
			}
			continue
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		size := int64(len(bodyBytes))
		bodyStr := string(bodyBytes)

		// Status kontrolü
		statusResult := assert.AssertStatus(tc.Expect.Status, resp.StatusCode)
		report.PrintResult(tc.Name, "status", statusResult, duration, size, tc.Expect.Status, resp.StatusCode)
		if !statusResult.Pass && config.FailFast {
			break abortedLoop
		}

		// Body kontrolü (varsa)
		if tc.Expect.Body != "" {
			bodyResult := assert.AssertBody(tc.Expect.Body, bodyStr)
			report.PrintResult(tc.Name, "body", bodyResult, duration, size, 0, 0)
			if !bodyResult.Pass && config.FailFast {
				break abortedLoop
			}
		}
	}
}

func RunAllTests(suite parser.TestSuite) {
	// 1) Run Test
	executeSuite(suite)

	// 2) reload
	reload := func() {
		report.ClearResults() // önceki sonuçları sil
		executeSuite(suite)   // aynı testleri tekrar çalıştır
	}

	// 3) Bubble Tea ekranını aç – reload’u geçir
	report.RenderResults(reload)
}
