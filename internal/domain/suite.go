package domain

type TestSuite struct {
	Name     string        `yaml:"suite_name"`
	Requests []TestRequest `yaml:"requests"`
}

type TestRequest struct {
	Name       string        `yaml:"name"`
	Req        HTTPRequest   `yaml:"request"`
	Expect     []Expectation `yaml:"expect"`
	DependsOn  []string      `yaml:"depends_on,omitempty"`
}

type HTTPRequest struct {
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Body    any               `yaml:"body,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type Expectation struct {
	Type   string                 `yaml:"expectation_type"`
	Kwargs map[string]any `yaml:"kwargs"`
}
