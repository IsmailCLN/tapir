package assert

import (
	"fmt"

	"github.com/IsmailCLN/tapir/internal/helpers"
)

func expectStatusCodeEquals(_ []byte, kw map[string]any) error {
	actual, ok := helpers.GetInt(kw, keyStatus)
	if !ok {
		return fmt.Errorf("expect_status_code_equals: %q was not injected or not an integer", keyStatus)
	}
	expected, ok := helpers.GetInt(kw, keyExpectedStatus)
	if !ok {
		return fmt.Errorf("expect_status_code_equals: missing or invalid %q", expected)
	}
	if actual != expected {
		return fmt.Errorf("status code mismatch: got=%d, want=%d", actual, expected)
	}
	return nil
}

func init() {
	Register("expect_status_code_equals", expectStatusCodeEquals)
}
