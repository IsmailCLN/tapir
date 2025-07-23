package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"github.com/IsmailCLN/tapir/models"
)

func RunTestCase(tc models.TestCase) {
	client := &http.Client{}

	var body io.Reader
	if tc.Body != "" {
		body = bytes.NewBuffer([]byte(tc.Body))
	}

	req, err := http.NewRequest(tc.Method, tc.URL, body)
	if err != nil {
		fmt.Printf("[FAIL] %s: request creation error: %v\n", tc.Name, err)
		return
	}

	for k, v := range tc.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[FAIL] %s: HTTP error: %v\n", tc.Name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == tc.Expect.Status {
		fmt.Printf("[PASS] %s\n", tc.Name)
	} else {
		fmt.Printf("[FAIL] %s: expected %d, got %d\n", tc.Name, tc.Expect.Status, resp.StatusCode)
	}
}
