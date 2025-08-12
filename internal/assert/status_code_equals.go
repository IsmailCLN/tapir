package assert

import (
	"fmt"

	"github.com/IsmailCLN/tapir/internal/helpers"
)

const (
	keyStatusToEqual = "status_code" // injected by runner
	keyExpected      = "code"        // YAML
)

func expectStatusCodeEquals(_ []byte, kw map[string]any) error {
	actual, ok := helpers.GetInt(kw, keyStatusToEqual)
	if !ok {
		return fmt.Errorf("expect_status_code_equals: %q was not injected or not an integer", keyStatusToEqual)
	}
	expected, ok := helpers.GetInt(kw, keyExpected)
	if !ok {
		return fmt.Errorf("expect_status_code_equals: missing or invalid %q", keyExpected)
	}
	if actual != expected {
		return fmt.Errorf("status code mismatch: got=%d, want=%d", actual, expected)
	}
	return nil
}

func init() {
	Register("expect_status_code_equals", expectStatusCodeEquals)
}
