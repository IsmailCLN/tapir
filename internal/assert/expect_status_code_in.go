package assert

import (
	"fmt"

	"github.com/IsmailCLN/tapir/internal/helpers"
)

const (
	keyCodes = "codes" // YAML: list or comma-separated string
)

func expectStatusCodeIn(_ []byte, kw map[string]any) error {
	code, ok := helpers.GetInt(kw, keyStatus)
	if !ok {
		return fmt.Errorf("expect_status_code_in: %q was not injected or not an integer", keyStatus)
	}

	allowed, ok := helpers.GetIntSlice(kw, keyCodes)
	if !ok || len(allowed) == 0 {
		return fmt.Errorf("expect_status_code_in: %q must be a non-empty list of integers", keyCodes)
	}

	for _, c := range allowed {
		if code == c {
			return nil
		}
	}
	return fmt.Errorf("status code %d is not in allowed set %v", code, allowed)
}

func init() { Register("expect_status_code_in", expectStatusCodeIn) }
