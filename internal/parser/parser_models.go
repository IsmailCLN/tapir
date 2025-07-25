package parser

type TestCase struct {
	Name    string            `yaml:"name"`
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
	Expect  Expectation       `yaml:"expect"`
}

type Expectation struct {
	Status int    `yaml:"status"`
	Body   string `yaml:"body,omitempty"`
}

type TestSuite struct {
	Tests []TestCase `yaml:"tests"`
}
