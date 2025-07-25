package assert

import (
	"strings"
)

type AssertionResult struct {
	Pass        bool
	Description string
	Expected    any
	Actual      any
}

func AssertBody(expected, actual string) AssertionResult {
	if expected == "" {
		return AssertionResult{Pass: true} // Body kontrolü istenmiyor
	}

	if strings.TrimSpace(expected) == strings.TrimSpace(actual) {
		return AssertionResult{Pass: true}
	}
	return AssertionResult{
		Pass:        false,
		Description: "body mismatch",
		Expected:    expected,
		Actual:      actual,
	}
}

func AssertStatus(expected, actual int) AssertionResult {
	if expected == actual {
		return AssertionResult{
			Pass:        true,
			Description: "Status code matches",
			Expected:    expected,
			Actual:      actual,
		}
	}

	return AssertionResult{
		Pass:        false,
		Description: "Status code does not match",
		Expected:    expected,
		Actual:      actual,
	}
}
