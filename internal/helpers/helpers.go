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

func String(m map[string]any, k string) (string, bool) {
	v, ok := m[k]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	return s, true
}

func Bool(kwargs map[string]any, key string) (bool, bool) {
	v, ok := kwargs[key]
	if !ok {
		return false, false
	}
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		s := strings.ToLower(strings.TrimSpace(t))
		switch s {
		case "true", "1", "yes", "y":
			return true, true
		case "false", "0", "no", "n":
			return false, true
		default:
			return false, false
		}
	default:
		return false, false
	}
}

func IntOpt(kwargs map[string]any, key string) (int, bool, error) {
	if _, ok := kwargs[key]; !ok {
		return 0, false, nil
	}
	v, err := Int(kwargs, key)
	if err != nil {
		return 0, true, err
	}
	return v, true, nil
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
