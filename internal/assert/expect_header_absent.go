package assert

import (
	"fmt"
	"strings"

	"github.com/IsmailCLN/tapir/internal/helpers"
)

func expectHeaderAbsent(_ []byte, kw map[string]any) error {
	hdr, err := helpers.AsHTTPHeader(kw[keyInjectedHeaders])
	if err != nil {
		return fmt.Errorf("expect_header_absent: response headers not available: %w", err)
	}
	name, ok := helpers.GetString(kw, "header")
	if !ok || strings.TrimSpace(name) == "" {
		return fmt.Errorf("expect_header_absent: missing or empty %q", "header")
	}

	if vals := hdr.Values(name); len(vals) > 0 {
		return fmt.Errorf("header %s should be absent, but present: %q", name, vals)
	}
	return nil
}

func init() { Register("expect_header_absent", expectHeaderAbsent) }