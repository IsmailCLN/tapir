package unit

import (
	"testing"

	"github.com/IsmailCLN/tapir/internal/assert"
)

func TestAssertStatus(t *testing.T) {
	cases := []struct {
		name     string
		exp, got int
		pass     bool
	}{
		{"match 200", 200, 200, true},
		{"mismatch 200/404", 200, 404, false},
		{"match 404", 404, 404, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := assert.AssertStatus(c.exp, c.got)
			if res.Pass != c.pass {
				t.Fatalf("expected pass=%v, got=%v (desc=%s)", c.pass, res.Pass, res.Description)
			}
		})
	}
}

func TestAssertBody(t *testing.T) {
	exp := `{"id":1,"name":"Leanne"}`
	got := `{"id":1,"name":"Leanne"}`
	if res := assert.AssertBody(exp, got); !res.Pass {
		t.Fatalf("should pass: %+v", res)
	}
	got2 := `{"id":2,"name":"Other"}`
	if res := assert.AssertBody(exp, got2); res.Pass {
		t.Fatalf("should fail")
	}
}
