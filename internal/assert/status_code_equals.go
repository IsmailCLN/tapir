package assert

import (
	"fmt"

	"github.com/IsmailCLN/tapir/internal/utils"
)

const (
	keyStatusToEqual = "status_code" // runner enjekte ediyor
	keyExpected      = "code"
)

func ExpectStatusCodeEquals(_ []byte, kw map[string]any) error {
	actual, err := utils.Int(kw, keyStatusToEqual)
	if err != nil {
		return err
	}
	expected, err := utils.Int(kw, keyExpected)
	if err != nil {
		return err
	}
	if actual != expected {
		return fmt.Errorf("status code %d ≠ %d", actual, expected)
	}
	return nil
}

func init() {
	Register("expect_status_code_equals", ExpectStatusCodeEquals)
}
