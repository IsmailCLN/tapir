package yaml

import (
	"os"

	"github.com/IsmailCLN/tapir/internal/domain"

	"gopkg.in/yaml.v2"
)

func LoadTestSuite(path string) (domain.TestSuite, error) {
	var suite domain.TestSuite

	data, err := os.ReadFile(path)
	if err != nil {
		return suite, err
	}

	err = yaml.Unmarshal(data, &suite)
	return suite, err
}
