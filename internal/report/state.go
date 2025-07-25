package report

import (
	"fmt"

	"github.com/IsmailCLN/tapir/internal/assert"
)

type TestResult struct {
	Name           string
	Category       string
	Result         assert.AssertionResult
	Duration       string
	ResponseSize   string
	ExpectedStatus int
	ActualStatus   int
}

var testResults []TestResult

func PrintResult(name, category string, result assert.AssertionResult, durationMs int64, sizeBytes int64, expected, actual int) {
	testResults = append(testResults, TestResult{
		Name:           name,
		Category:       category,
		Result:         result,
		Duration:       formatDuration(durationMs),
		ResponseSize:   formatSize(sizeBytes),
		ExpectedStatus: expected,
		ActualStatus:   actual,
	})
}

func ClearResults() {
	testResults = nil
}

func GetResults() []TestResult {
	return testResults
}

func formatDuration(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%d ms", ms)
	}
	return fmt.Sprintf("%.2f s", float64(ms)/1000.0)
}

func formatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	return fmt.Sprintf("%.2f KB", float64(bytes)/1024.0)
}
