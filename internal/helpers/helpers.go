package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

// Int returns kwargs[key] as int, or an explanatory error.
func Int(kwargs map[string]any, key string) (int, error) {
	v, ok := kwargs[key]
	if !ok {
		return 0, fmt.Errorf("missing parameter: %s", key)
	}
	switch t := v.(type) {
	case int:
		return t, nil
	case int64:
		return int(t), nil
	case float64:
		return int(t), nil
	case string:
		i, err := strconv.Atoi(t)
		if err != nil {
			return 0, fmt.Errorf("%s must be int, got %q", key, t)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("%s must be int, got %T", key, v)
	}
}

// sanitize: makes \n, \r\n, \t ve \" more readable.
func Sanitize(s string) string {
	replacer := strings.NewReplacer(
		`\"`, `"`,
		`\r\n`, "",
		`\n`, "",
		`\t`, "",
	)
	return replacer.Replace(s)
}
