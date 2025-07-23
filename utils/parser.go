package utils

import (
	"os"

	"github.com/IsmailCLN/tapir/models"
	"gopkg.in/yaml.v2"
)

func LoadTestSuite(path string) (models.TestSuite, error) {
	var suite models.TestSuite

	data, err := os.ReadFile(path)
	if err != nil {
		return suite, err
	}

	err = yaml.Unmarshal(data, &suite)
	return suite, err
}
