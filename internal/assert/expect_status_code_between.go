package assert

import (
	"fmt"

	"github.com/IsmailCLN/tapir/internal/helpers"
)

func ExpectStatusCodeBetween(_ []byte, kw map[string]any) error {
	code, ok := helpers.GetInt(kw, keyStatus)
	if !ok {
		return fmt.Errorf("expect_status_code_between: %q was not injected or not an integer", keyStatus)
	}
	min, ok := helpers.GetInt(kw, keyMin)
	if !ok {
		return fmt.Errorf("expect_status_code_between: missing or invalid %q", keyMin)
	}
	max, ok := helpers.GetInt(kw, keyMax)
	if !ok {
		return fmt.Errorf("expect_status_code_between: missing or invalid %q", keyMax)
	}
	if min > max {
		return fmt.Errorf("expect_status_code_between: %q must be <= %q (got %d > %d)", keyMin, keyMax, min, max)
	}

	if code < min || code > max {
		return fmt.Errorf("status %d is not within [%d..%d]", code, min, max)
	}
	return nil
}

func init() {
	Register("expect_status_code_between", ExpectStatusCodeBetween)
}
