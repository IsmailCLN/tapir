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

func getCookieName(kwargs map[string]any) (string, error) {
	keys := []string{"cookieName", "cookie_name"}
	for _, k := range keys {
		if v, ok := kwargs[k]; ok {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				return s, nil
			}
		}
	}
	return "", fmt.Errorf("missing parameter %q (string)", "cookieName")
}

func coerceHeaders(v any) (http.Header, error) {
	switch h := v.(type) {
	case http.Header:
		return h, nil
	case map[string][]string:
		return http.Header(h), nil
	case map[string]any:
		out := make(http.Header, len(h))
		for k, vv := range h {
			switch t := vv.(type) {
			case []string:
				out[k] = t
			case string:
				out[k] = []string{t}
			}
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported headers type: %T", v)
	}
}
