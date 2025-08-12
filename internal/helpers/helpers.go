package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
Coercion tier:
  - AsX(v any) (T, error): convert a loose-typed value to a target type.
  - GetX(kwargs map[string]any, key string) (T, bool): read from kwargs and coerce.
Compatibility shims (older call sites):
  - String/Bool/Int/... mirror GetX signature.
*/

// ---------- Atomics ----------

func AsString(v any) (string, error) {
	switch t := v.(type) {
	case string:
		return t, nil
	case []byte:
		return string(t), nil
	case fmt.Stringer:
		return t.String(), nil
	case int, int32, int64, float32, float64, bool, uint, uint32, uint64:
		return fmt.Sprint(t), nil
	case nil:
		return "", errors.New("string: value is nil")
	default:
		return "", fmt.Errorf("string: unsupported type %T", v)
	}
}

func AsInt(v any) (int, error) {
	switch t := v.(type) {
	case int:
		return t, nil
	case int64:
		return int(t), nil
	case int32:
		return int(t), nil
	case float64:
		return int(t), nil
	case float32:
		return int(t), nil
	case uint:
		return int(t), nil
	case uint32:
		return int(t), nil
	case uint64:
		return int(t), nil
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return 0, errors.New("int: empty string")
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("int: parse error for %q: %w", s, err)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("int: unsupported type %T", v)
	}
}

func AsFloat64(v any) (float64, error) {
	switch t := v.(type) {
	case float64:
		return t, nil
	case float32:
		return float64(t), nil
	case int:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case uint, uint32, uint64:
		return AsFloat64(fmt.Sprint(t))
	case string:
		s := strings.TrimSpace(t)
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("float: parse error for %q: %w", s, err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("float: unsupported type %T", v)
	}
}

func AsBool(v any) (bool, error) {
	switch t := v.(type) {
	case bool:
		return t, nil
	case string:
		s := strings.ToLower(strings.TrimSpace(t))
		switch s {
		case "true", "t", "yes", "y", "on", "1":
			return true, nil
		case "false", "f", "no", "n", "off", "0", "":
			return false, nil
		default:
			return false, fmt.Errorf("bool: invalid string %q", t)
		}
	case int, int32, int64:
		i, _ := AsInt(t)
		return i != 0, nil
	case float32, float64:
		f, _ := AsFloat64(t)
		return f != 0, nil
	default:
		return false, fmt.Errorf("bool: unsupported type %T", v)
	}
}

// Duration parsing: strings via time.ParseDuration, numeric as milliseconds.
func AsDuration(v any) (time.Duration, error) {
	switch t := v.(type) {
	case time.Duration:
		return t, nil
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return 0, errors.New("duration: empty string")
		}
		d, err := time.ParseDuration(s)
		if err != nil {
			return 0, fmt.Errorf("duration: parse error for %q: %w", s, err)
		}
		return d, nil
	case int, int64, float64, float32:
		ms, _ := AsFloat64(t)
		return time.Duration(ms) * time.Millisecond, nil
	default:
		return 0, fmt.Errorf("duration: unsupported type %T", v)
	}
}

// ---------- Slices ----------

func AsStringSlice(v any) ([]string, error) {
	switch t := v.(type) {
	case []string:
		return t, nil
	case []any:
		out := make([]string, 0, len(t))
		for _, e := range t {
			s, err := AsString(e)
			if err != nil {
				return nil, err
			}
			out = append(out, s)
		}
		return out, nil
	case []float64:
		out := make([]string, 0, len(t))
		for _, f := range t {
			out = append(out, fmt.Sprint(f))
		}
		return out, nil
	case string:
		// CSV-ish: "a,b, c"
		parts := strings.Split(t, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
		if len(out) == 0 {
			return nil, errors.New("string slice: empty after split")
		}
		return out, nil
	default:
		return nil, fmt.Errorf("string slice: unsupported type %T", v)
	}
}

func AsIntSlice(v any) ([]int, error) {
	switch t := v.(type) {
	case []int:
		return t, nil
	case []any:
		out := make([]int, 0, len(t))
		for _, e := range t {
			i, err := AsInt(e)
			if err != nil {
				return nil, err
			}
			out = append(out, i)
		}
		return out, nil
	case []float64:
		out := make([]int, 0, len(t))
		for _, f := range t {
			out = append(out, int(f))
		}
		return out, nil
	case []string:
		out := make([]int, 0, len(t))
		for _, s := range t {
			i, err := AsInt(s)
			if err != nil {
				return nil, err
			}
			out = append(out, i)
		}
		return out, nil
	case string:
		parts := strings.Split(t, ",")
		out := make([]int, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			i, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			out = append(out, i)
		}
		if len(out) == 0 {
			return nil, errors.New("int slice: empty after split")
		}
		return out, nil
	default:
		return nil, fmt.Errorf("int slice: unsupported type %T", v)
	}
}

// ---------- Maps / Headers ----------

func AsMapStringAny(v any) (map[string]any, error) {
	switch t := v.(type) {
	case map[string]any:
		return t, nil
	case map[any]any:
		m := make(map[string]any, len(t))
		for k, v := range t {
			ks, err := AsString(k)
			if err != nil {
				return nil, fmt.Errorf("map key not stringable: %v", k)
			}
			m[ks] = v
		}
		return m, nil
	default:
		return nil, fmt.Errorf("map[string]any: unsupported type %T", v)
	}
}

// Accept http.Header, map[string]string, map[string][]string, map[string]any (string or []string)
func AsHTTPHeader(v any) (http.Header, error) {
	switch t := v.(type) {
	case http.Header:
		return t, nil
	case map[string][]string:
		h := http.Header{}
		for k, vals := range t {
			for _, val := range vals {
				h.Add(k, val)
			}
		}
		return h, nil
	case map[string]string:
		h := http.Header{}
		for k, val := range t {
			h.Add(k, val)
		}
		return h, nil
	case map[string]any:
		h := http.Header{}
		for k, val := range t {
			switch vv := val.(type) {
			case string:
				h.Add(k, vv)
			case []string:
				for _, s := range vv {
					h.Add(k, s)
				}
			case []any:
				for _, e := range vv {
					s, err := AsString(e)
					if err != nil {
						return nil, fmt.Errorf("header %s value not stringable: %v", k, e)
					}
					h.Add(k, s)
				}
			default:
				s, err := AsString(vv)
				if err != nil {
					return nil, fmt.Errorf("header %s value not stringable: %T", k, vv)
				}
				h.Add(k, s)
			}
		}
		return h, nil
	default:
		return nil, fmt.Errorf("http.Header: unsupported type %T", v)
	}
}

// ---------- Dotted path lookup (maps/slices) ----------

// LookupPath navigates a JSON-like structure using dotted path (e.g., "data.items.0.id").
// Supports map[string]any and []any; numeric segments index slices.
func LookupPath(root any, path string) (any, bool) {
	if path == "" {
		return root, true
	}
	cur := root
	for _, seg := range strings.Split(path, ".") {
		if m, ok := cur.(map[string]any); ok {
			nxt, exists := m[seg]
			if !exists {
				return nil, false
			}
			cur = nxt
			continue
		}
		if arr, ok := cur.([]any); ok {
			// seg must be index
			idx, err := strconv.Atoi(seg)
			if err != nil || idx < 0 || idx >= len(arr) {
				return nil, false
			}
			cur = arr[idx]
			continue
		}
		// tolerate map[any]any
		if mm, ok := cur.(map[any]any); ok {
			var hit any
			found := false
			for k, v := range mm {
				ks, _ := AsString(k)
				if ks == seg {
					hit, found = v, true
					break
				}
			}
			if !found {
				return nil, false
			}
			cur = hit
			continue
		}
		return nil, false
	}
	return cur, true
}

// ---------- GetX helpers over kwargs ----------

func GetString(kwargs map[string]any, key string) (string, bool) {
	v, ok := kwargs[key]
	if !ok {
		return "", false
	}
	s, err := AsString(v)
	return s, err == nil
}

func GetInt(kwargs map[string]any, key string) (int, bool) {
	v, ok := kwargs[key]
	if !ok {
		return 0, false
	}
	i, err := AsInt(v)
	return i, err == nil
}

func GetFloat64(kwargs map[string]any, key string) (float64, bool) {
	v, ok := kwargs[key]
	if !ok {
		return 0, false
	}
	f, err := AsFloat64(v)
	return f, err == nil
}

func GetBool(kwargs map[string]any, key string) (bool, bool) {
	v, ok := kwargs[key]
	if !ok {
		return false, false
	}
	b, err := AsBool(v)
	return b, err == nil
}

func GetStringSlice(kwargs map[string]any, key string) ([]string, bool) {
	v, ok := kwargs[key]
	if !ok {
		return nil, false
	}
	s, err := AsStringSlice(v)
	return s, err == nil
}

func GetIntSlice(kwargs map[string]any, key string) ([]int, bool) {
	v, ok := kwargs[key]
	if !ok {
		return nil, false
	}
	s, err := AsIntSlice(v)
	return s, err == nil
}

func GetDuration(kwargs map[string]any, key string) (time.Duration, bool) {
	v, ok := kwargs[key]
	if !ok {
		return 0, false
	}
	d, err := AsDuration(v)
	return d, err == nil
}

// sanitize: makes \n, \r\n, \t ve \" more readable.
func Sanitize(s string) string {
	replacer := strings.NewReplacer(
		`\"`, `"`,
		`\r\n`, "",
		`\n`, "",
		`\t`, "",
	)
	return replacer.Replace(s)
}
