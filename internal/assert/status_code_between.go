package assert

import (
	"fmt"

	"github.com/IsmailCLN/tapir/internal/helpers"
)

const (
	keyStatus = "status_code" // runner enjekte ediyor
	keyMin    = "min"         // YAML
	keyMax    = "max"         // YAML
)

func ExpectStatusCodeBetween(_ []byte, kw map[string]any) error {
	code, err := helpers.Int(kw, keyStatus)
	if err != nil {
		return err
	}
	min, err := helpers.Int(kw, keyMin)
	if err != nil {
		return err
	}
	max, err := helpers.Int(kw, keyMax)
	if err != nil {
		return err
	}

	if code < min || code > max {
		return fmt.Errorf("status %d ∉ [%d–%d]", code, min, max)
	}
	return nil
}

func init() {
	Register("expect_status_code_between", ExpectStatusCodeBetween)
}
