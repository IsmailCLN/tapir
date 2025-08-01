package assert

import "fmt"

const (
	keyStatus = "status_code" // runner enjekte ediyor
	keyMin    = "min"         // YAML
	keyMax    = "max"         // YAML
)

func ExpectStatusCodeBetween(_ []byte, kwargs map[string]interface{}) error {
	getInt := func(k string) (int, error) {
		v, ok := kwargs[k]
		if !ok {
			return 0, fmt.Errorf("parametre eksik: %s", k)
		}
		i, ok := v.(int)
		if !ok {
			return 0, fmt.Errorf("%s parametresi int değil", k)
		}
		return i, nil
	}

	code, err := getInt(keyStatus)
	if err != nil {
		return err
	}
	min, err := getInt(keyMin)
	if err != nil {
		return err
	}
	max, err := getInt(keyMax)
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
