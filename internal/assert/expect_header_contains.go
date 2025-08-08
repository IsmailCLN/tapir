package assert

import (
	"fmt"
	"net/http"
	"strings"
)

// Uses the same kw keys as expect_header_equals:
//   headers (injected by runner), header, value, ignore_case (optional)

// expect_header_contains: checks if a response header contains a substring.
// Kwargs:
//   header: string (required) -> e.g., "Cache-Control"
//   value:  string (required) -> substring to look for
//   ignore_case: bool (optional, default: false)
func expectHeaderContains(_ []byte, kw map[string]any) error {
	hdrs, ok := kw[keyInjectedHeaders].(http.Header)
	if !ok {
		return fmt.Errorf("expect_header_contains: response headers were not injected by the runner")
	}

	name, ok := kw[keyHeaderName].(string)
	if !ok || strings.TrimSpace(name) == "" {
		return fmt.Errorf("expect_header_contains: missing or empty %q parameter", keyHeaderName)
	}

	needle, ok := kw[keyExpectedValue].(string)
	if !ok || needle == "" {
		return fmt.Errorf("expect_header_contains: %q parameter must be a non-empty string", keyExpectedValue)
	}

	ignoreCase := false
	if v, ok := kw["ignore_case"].(bool); ok {
		ignoreCase = v
	}

	values := hdrs.Values(name)
	if len(values) == 0 {
		return fmt.Errorf("header %s not found", name)
	}

	if ignoreCase {
		needle = strings.ToLower(needle)
		for _, v := range values {
			if strings.Contains(strings.ToLower(v), needle) {
				return nil
			}
		}
	} else {
		for _, v := range values {
			if strings.Contains(v, needle) {
				return nil
			}
		}
	}

	return fmt.Errorf("header %s does not contain %q; got: %q", name, needle, strings.Join(values, ", "))
}

func init() { Register("expect_header_contains", expectHeaderContains) }
