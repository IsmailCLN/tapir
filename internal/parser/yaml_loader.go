package parser

import (
	"os"

	"github.com/IsmailCLN/tapir/internal/domain"
	"gopkg.in/yaml.v3"
)

func LoadTestSuite(path string) ([]domain.TestSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var suites []domain.TestSuite
	if err := yaml.Unmarshal(data, &suites); err != nil {
		return nil, err
	}
	return suites, nil
}
