package assert

import (
	"fmt"

	"github.com/IsmailCLN/tapir/internal/utils"
)

const (
	keyStatus = "status_code" // runner enjekte ediyor
	keyMin    = "min"         // YAML
	keyMax    = "max"         // YAML
)

func ExpectStatusCodeBetween(_ []byte, kw map[string]any) error {
	code, err := utils.Int(kw, keyStatus)
	if err != nil {
		return err
	}
	min, err := utils.Int(kw, keyMin)
	if err != nil {
		return err
	}
	max, err := utils.Int(kw, keyMax)
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
