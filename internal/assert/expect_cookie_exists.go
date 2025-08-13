package assert

import (
	"fmt"
	"net/http"
	"strings"
)

func init() {
	Register("expect_cookie_exists", expectCookieExists)
}

func expectCookieExists(_ []byte, kwargs map[string]any) error {
	// cookieName / cookie_name
	cname, err := getCookieName(kwargs)
	if err != nil {
		return err
	}

	// optional flag
	ignoreCase := false
	if v, ok := kwargs["ignore_case"]; ok {
		if b, ok := v.(bool); ok {
			ignoreCase = b
		}
	}

	h, err := coerceHeaders(kwargs["headers"])
	if err != nil {
		return fmt.Errorf("headers not available for cookie assertion: %w", err)
	}

	resp := &http.Response{Header: h}
	cookies := resp.Cookies()
	if ignoreCase {
		target := strings.ToLower(cname)
		for _, c := range cookies {
			if strings.ToLower(c.Name) == target {
				return nil
			}
		}
	} else {
		for _, c := range cookies {
			if c.Name == cname {
				return nil
			}
		}
	}

	return fmt.Errorf("expected cookie %q to exist but it was not found", cname)
}


