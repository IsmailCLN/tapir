package assert

import (
	"fmt"
	"net/http"
	"strings"
)

func init() { Register("expect_cookie_not_exists", expectCookieNotExists) }

func expectCookieNotExists(_ []byte, kwargs map[string]any) error {
	name, err := getCookieName(kwargs)
	if err != nil {
		return err
	}
	h, err := coerceHeaders(kwargs["headers"])
	if err != nil {
		return fmt.Errorf("headers not available: %w", err)
	}

	resp := &http.Response{Header: h}
	for _, c := range resp.Cookies() {
		if strings.EqualFold(c.Name, name) {
			return fmt.Errorf("cookie %q should not exist, but was found", name)
		}
	}
	return nil
}
