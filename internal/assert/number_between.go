package assert

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/IsmailCLN/tapir/internal/helpers"
)

func numberBetween(body []byte, kwargs map[string]any) error {
	field, ok := helpers.GetString(kwargs, "column")
	if !ok || field == "" {
		return fmt.Errorf("expect_number_to_be_between: missing or empty %q", "column")
	}
	min, ok := helpers.GetFloat64(kwargs, "min")
	if !ok {
		return fmt.Errorf("expect_number_to_be_between: missing or invalid %q", "min")
	}
	max, hasMax := helpers.GetFloat64(kwargs, "max")

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		return fmt.Errorf("expect_number_to_be_between: invalid JSON: %v", err)
	}

	raw, exists := m[field]
	if !exists {
		return fmt.Errorf("field %s not found or not numeric", field)
	}
	val, err := helpers.AsFloat64(raw)
	if err != nil {
		if s, e2 := helpers.AsString(raw); e2 == nil {
			if v2, e3 := helpers.AsFloat64(s); e3 == nil {
				val = v2
			} else {
				return fmt.Errorf("field %s not found or not numeric", field)
			}
		} else {
			return fmt.Errorf("field %s not found or not numeric", field)
		}
	}

	outside := val < min || (hasMax && val > max)
	if outside {
		bounds := fmt.Sprintf("[%.6g, +inf)", min)
		if hasMax {
			bounds = fmt.Sprintf("[%.6g, %.6g]", min, max)
		}
		return fmt.Errorf("%s=%.6g outside %s", field, val, bounds)
	}
	return nil
}

func init() { Register("expect_number_to_be_between", numberBetween) }
