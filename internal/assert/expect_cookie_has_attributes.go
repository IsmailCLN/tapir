package assert

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IsmailCLN/tapir/internal/helpers"
)

func init() { Register("expect_cookie_has_attributes", expectCookieHasAttributes) }

func expectCookieHasAttributes(_ []byte, kwargs map[string]any) error {
	name, err := getCookieNameCompat(kwargs)
	if err != nil {
		return err
	}

	hdr, err := helpers.AsHTTPHeader(kwargs["headers"])
	if err != nil {
		return fmt.Errorf("headers not available: %w", err)
	}

	wantPath, _ := helpers.GetString(kwargs, "path")
	wantDomain, _ := helpers.GetString(kwargs, "domain")
	wantHttpOnly, hasHttpOnly := helpers.GetBool(kwargs, "http_only")
	wantSecure, hasSecure := helpers.GetBool(kwargs, "secure")
	wantSameSiteStr, hasSameSite := helpers.GetString(kwargs, "samesite")

	// min_max_age: distinguish "absent" vs "present but invalid"
	var (
		minMaxAge     int
		hasMinMaxAge  bool
	)
	if _, present := kwargs["min_max_age"]; present {
		v, err := helpers.AsInt(kwargs["min_max_age"])
		if err != nil {
			return fmt.Errorf("min_max_age: %v", err)
		}
		minMaxAge = v
		hasMinMaxAge = true
	}

	requireNotExpired, hasRequireNotExpired := helpers.GetBool(kwargs, "not_expired")

	var wantSameSite http.SameSite
	if hasSameSite {
		ss, err := parseSameSite(wantSameSiteStr)
		if err != nil {
			return err
		}
		wantSameSite = ss
	}

	now := time.Now()
	resp := &http.Response{Header: hdr}

	for _, c := range resp.Cookies() {
		if c.Name != name {
			continue
		}

		if wantPath != "" && c.Path != wantPath {
			return fmt.Errorf("cookie %q path mismatch: got %q want %q", name, c.Path, wantPath)
		}
		if wantDomain != "" && !strings.EqualFold(c.Domain, wantDomain) {
			return fmt.Errorf("cookie %q domain mismatch: got %q want %q", name, c.Domain, wantDomain)
		}
		if hasHttpOnly && c.HttpOnly != wantHttpOnly {
			return fmt.Errorf("cookie %q httponly mismatch: got %v want %v", name, c.HttpOnly, wantHttpOnly)
		}
		if hasSecure && c.Secure != wantSecure {
			return fmt.Errorf("cookie %q secure mismatch: got %v want %v", name, c.Secure, wantSecure)
		}
		if hasSameSite && c.SameSite != wantSameSite {
			return fmt.Errorf("cookie %q samesite mismatch: got %s want %s",
				name, sameSiteString(c.SameSite), sameSiteString(wantSameSite))
		}

		if hasMinMaxAge {
			switch {
			case c.MaxAge > 0:
				if c.MaxAge < minMaxAge {
					return fmt.Errorf("cookie %q max-age too small: got %d want >= %d", name, c.MaxAge, minMaxAge)
				}
			case !c.Expires.IsZero():
				remaining := int(time.Until(c.Expires).Seconds())
				if remaining < minMaxAge {
					return fmt.Errorf("cookie %q expires too soon: remaining %ds want >= %ds", name, remaining, minMaxAge)
				}
			default:
				return fmt.Errorf("cookie %q has neither Max-Age nor Expires to validate min_max_age", name)
			}
		}

		if hasRequireNotExpired {
			expired := isCookieExpired(c, now)
			if requireNotExpired && expired {
				return fmt.Errorf("cookie %q already expired: max-age=%d expires=%s", name, c.MaxAge, c.Expires)
			}
			if !requireNotExpired && !expired {
				return fmt.Errorf("cookie %q is not expired", name)
			}
		}

		return nil
	}

	return fmt.Errorf("cookie %q not found", name)
}

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
