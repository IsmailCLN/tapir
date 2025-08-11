package assert

import (
	"fmt"
	"net/http"
	"strings"
)

func init() { Register("expect_cookie_value_equals", expectCookieValueEquals) }

func expectCookieValueEquals(_ []byte, kwargs map[string]any) error {
	name, err := getCookieName(kwargs)
	if err != nil {
		return err
	}

	want, ok := kwargs["value"].(string)
	if !ok || strings.TrimSpace(want) == "" {
		return fmt.Errorf("missing parameter %q (string)", "value")
	}

	h, err := coerceHeaders(kwargs["headers"])
	if err != nil {
		return fmt.Errorf("headers not available: %w", err)
	}

	resp := &http.Response{Header: h}
	for _, c := range resp.Cookies() {
		if c.Name == name {
			if c.Value == want {
				return nil
			}
			return fmt.Errorf("cookie %q value mismatch: got %q want %q", name, c.Value, want)
		}
	}
	return fmt.Errorf("cookie %q not found", name)
}
