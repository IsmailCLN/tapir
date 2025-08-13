package assert

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IsmailCLN/tapir/internal/helpers"
)

func getCookieNameCompat(kwargs map[string]any) (string, error) {
	if s, ok := helpers.GetString(kwargs, "cookieName"); ok && strings.TrimSpace(s) != "" {
		return s, nil
	}
	if s, ok := helpers.GetString(kwargs, "name"); ok && strings.TrimSpace(s) != "" {
		return s, nil
	}
	return "", fmt.Errorf("expect_cookie_has_attributes: missing or empty %q", "cookieName")
}

func isCookieExpired(c *http.Cookie, now time.Time) bool {
	if c.MaxAge < 0 {
		return true
	}
	if !c.Expires.IsZero() && !c.Expires.After(now) {
		return true
	}
	return false
}

func parseSameSite(s string) (http.SameSite, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "default", "defaultmode":
		return http.SameSiteDefaultMode, nil
	case "lax", "laxmode":
		return http.SameSiteLaxMode, nil
	case "strict", "strictmode":
		return http.SameSiteStrictMode, nil
	case "none", "nonemode":
		return http.SameSiteNoneMode, nil
	default:
		return 0, fmt.Errorf("invalid samesite value %q (use one of: default|lax|strict|none)", s)
	}
}

func sameSiteString(ss http.SameSite) string {
	switch ss {
	case http.SameSiteDefaultMode:
		return "Default"
	case http.SameSiteLaxMode:
		return "Lax"
	case http.SameSiteStrictMode:
		return "Strict"
	case http.SameSiteNoneMode:
		return "None"
	default:
		return fmt.Sprintf("SameSite(%d)", int(ss))
	}
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
