package assert

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	keyInjectedHeaders = "headers"
	keyHeaderName      = "header"
	keyExpectedValue   = "value"
)

// expect_header_equals: checks if a response header equals the expected value.
// Kwargs:
//
//	header: string (required) -> e.g., "Content-Type"
//	value:  string (required)
//	ignore_case: bool (optional, default: false)
func expectHeaderEquals(_ []byte, kw map[string]any) error {
	hdrs, ok := kw[keyInjectedHeaders].(http.Header)
	if !ok {
		return fmt.Errorf("expect_header_equals: response headers were not injected by the runner")
	}

	name, ok := kw[keyHeaderName].(string)
	if !ok || name == "" {
		return fmt.Errorf("expect_header_equals: missing or empty %q parameter", keyHeaderName)
	}

	want, ok := kw[keyExpectedValue].(string)
	if !ok {
		return fmt.Errorf("expect_header_equals: %q parameter must be a string", keyExpectedValue)
	}

	ignoreCase := false
	if v, ok := kw["ignore_case"].(bool); ok {
		ignoreCase = v
	}

	got := hdrs.Get(name) // case-insensitive lookup
	if ignoreCase {
		if !strings.EqualFold(got, want) {
			return fmt.Errorf("header %s mismatch (case-insensitive): got=%q, want=%q", name, got, want)
		}
	} else {
		if got != want {
			return fmt.Errorf("header %s mismatch: got=%q, want=%q", name, got, want)
		}
	}
	return nil
}

func init() { Register("expect_header_equals", expectHeaderEquals) }
