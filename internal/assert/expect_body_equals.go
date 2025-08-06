package assert

import (
	"fmt"
	"regexp"
)

var spaceRE = regexp.MustCompile(`\s+`)

func expectBodyEquals(body []byte, kw map[string]any) error {
	v, ok := kw["value"]
	if !ok {
		return fmt.Errorf("expect_body_equals: 'value' parametresi eksik")
	}

	want, ok := v.(string)
	if !ok {
		return fmt.Errorf("expect_body_equals: 'value' parametresi string olmalı")
	}

	clean := func(s string) string { return spaceRE.ReplaceAllString(s, "") }
	if clean(string(body)) != clean(want) {
		return fmt.Errorf("body eşleşmedi:\nwant=%q\ngot=%q", want, string(body))
	}
	return nil
}

func init() { Register("expect_body_equals", expectBodyEquals) }
