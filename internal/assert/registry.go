package assert

type Func func(respBody []byte, kwargs map[string]interface{}) error

var registry = map[string]Func{}

func Register(name string, f Func) { registry[name] = f }

func Get(name string) (Func, bool) { f, ok := registry[name]; return f, ok }
