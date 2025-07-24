package assert

type AssertionResult struct {
	Pass        bool
	Description string
	Expected    any
	Actual      any
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
