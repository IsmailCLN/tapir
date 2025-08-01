package assert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

func init() { Register("expect_number_to_be_between", numberBetween) }

func toFloat(v interface{}) (float64, bool) {
	switch t := v.(type) {

	// sayısal tipler
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case float32:
		return float64(t), true
	case float64:
		return t, true

	case json.Number:
		f, err := t.Float64()
		return f, err == nil

	case string:
		f, err := strconv.ParseFloat(t, 64)
		return f, err == nil

	default:
		return math.NaN(), false
	}
}

func numberBetween(body []byte, kwargs map[string]interface{}) error {
	field, _ := kwargs["column"].(string)
	min, _ := toFloat(kwargs["min"])
	max, maxOk := toFloat(kwargs["max"])

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()

	var m map[string]interface{}
	if err := dec.Decode(&m); err != nil {
		return err
	}

	val, ok := toFloat(m[field])
	if !ok {
		return fmt.Errorf("field %s not found or not numeric", field)
	}

	if val < min || (maxOk && val > max) {
		return fmt.Errorf("%s=%v outside [%v,%v]", field, val, min, max)
	}
	return nil
}
