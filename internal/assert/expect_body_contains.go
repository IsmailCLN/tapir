package assert

import (
	"fmt"
	"strings"
)

func expectBodyContains(body []byte, kw map[string]any) error {
	v, ok := kw["value"]
	if !ok {
		return fmt.Errorf("expect_body_contains: 'value' parametresi eksik")
	}
	substr, ok := v.(string)
	if !ok {
		return fmt.Errorf("expect_body_contains: 'value' parametresi string olmalı")
	}

	clean := func(s string) string { return spaceRE.ReplaceAllString(s, "") }
	if !strings.Contains(clean(string(body)), clean(substr)) {
		return fmt.Errorf("body içinde %q bulunamadı", substr)
	}
	return nil
}

func init() { Register("expect_body_contains", expectBodyContains) }
