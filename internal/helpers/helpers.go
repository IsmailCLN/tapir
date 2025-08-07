package helpers

import (
	"fmt"
	"strings"
)

// Int returns kwargs[key] as int, or an explanatory error.
func Int(kwargs map[string]any, key string) (int, error) {
	v, ok := kwargs[key]
	if !ok {
		return 0, fmt.Errorf("missing parameter: %s", key)
	}
	i, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("%s parameter is not an int", key)
	}
	return i, nil
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
