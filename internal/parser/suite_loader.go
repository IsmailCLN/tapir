package parser

import (
	"os"

	"gopkg.in/yaml.v2"
)

func LoadTestSuite(path string) (TestSuite, error) {
	var suite TestSuite

	data, err := os.ReadFile(path)
	if err != nil {
		return suite, err
	}

	err = yaml.Unmarshal(data, &suite)
	return suite, err
}
