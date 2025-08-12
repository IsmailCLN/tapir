package assert

import (
	"fmt"
	"mime"
	"strings"

	"github.com/IsmailCLN/tapir/internal/helpers"
)

func expectContentType(_ []byte, kw map[string]any) error {
	hdr, err := helpers.AsHTTPHeader(kw[keyInjectedHeaders])
	if err != nil {
		return fmt.Errorf("expect_content_type: response headers not available: %w", err)
	}

	want, ok := helpers.GetString(kw, "value")
	if !ok || strings.TrimSpace(want) == "" {
		return fmt.Errorf("expect_content_type: missing or empty %q", "value")
	}
	ignoreParams, _ := helpers.GetBool(kw, "ignore_params") // default false
	ignoreCase, _ := helpers.GetBool(kw, "ignore_case")     // default false

	got := hdr.Get("Content-Type")
	if got == "" {
		return fmt.Errorf("Content-Type header not found")
	}

	if ignoreParams {
		mtGot, _, err := mime.ParseMediaType(got)
		if err == nil {
			got = mtGot
		}
		mtWant, _, err := mime.ParseMediaType(want)
		if err == nil {
			want = mtWant
		}
	}

	if ignoreCase {
		if !strings.EqualFold(got, want) {
			return fmt.Errorf("content type mismatch (ci): got=%q, want=%q", got, want)
		}
	} else if got != want {
		return fmt.Errorf("content type mismatch: got=%q, want=%q", got, want)
	}
	return nil
}

func init() { Register("expect_content_type", expectContentType) }
