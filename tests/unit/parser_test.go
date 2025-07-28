package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IsmailCLN/tapir/internal/parser"
)

const goodYAML = `tests:
  - name: Get Users
    method: GET
    url: https://example.test/users
    expect:
      status: 200
`

func TestParser_Load_OK(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "good.yaml")
	if err := os.WriteFile(p, []byte(goodYAML), 0644); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	suite, err := parser.LoadTestSuite(p)
	if err != nil {
		t.Fatalf("LoadTestSuite failed: %v", err)
	}
	if len(suite.Tests) != 1 {
		t.Fatalf("expected 1 test, got %d", len(suite.Tests))
	}
	if suite.Tests[0].Expect.Status != 200 {
		t.Fatalf("expected status 200, got %d", suite.Tests[0].Expect.Status)
	}
}
