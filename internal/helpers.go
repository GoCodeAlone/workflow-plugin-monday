package internal

import "fmt"

func getModuleName(config map[string]any) string {
	if v, ok := config["module"].(string); ok && v != "" {
		return v
	}
	return "monday"
}

func resolveValue(key string, current, config map[string]any) string {
	if v, ok := current[key].(string); ok && v != "" {
		return v
	}
	if v, ok := config[key].(string); ok && v != "" {
		return v
	}
	return ""
}

func resolveInt64(key string, current, config map[string]any) int64 {
	if v := toInt64(current[key]); v != 0 {
		return v
	}
	return toInt64(config[key])
}

func resolveInt(key string, current, config map[string]any) int {
	return int(resolveInt64(key, current, config))
}

func resolveFloat64(key string, current, config map[string]any) float64 {
	if v := toFloat64(current[key]); v != 0 {
		return v
	}
	return toFloat64(config[key])
}

func resolveBool(key string, current, config map[string]any) bool {
	for _, m := range []map[string]any{current, config} {
		if v, ok := m[key]; ok {
			switch t := v.(type) {
			case bool:
				return t
			case string:
				return t == "true" || t == "1" || t == "yes"
			}
		}
	}
	return false
}

func resolveStringSlice(key string, current, config map[string]any) []string {
	for _, m := range []map[string]any{current, config} {
		if v, ok := m[key]; ok {
			switch t := v.(type) {
			case []string:
				return t
			case []any:
				result := make([]string, 0, len(t))
				for _, item := range t {
					if s, ok := item.(string); ok {
						result = append(result, s)
					}
				}
				return result
			}
		}
	}
	return nil
}

func resolveMap(key string, current, config map[string]any) map[string]any {
	for _, m := range []map[string]any{current, config} {
		if v, ok := m[key]; ok {
			if result, ok := v.(map[string]any); ok {
				return result
			}
		}
	}
	return nil
}

func toInt64(v any) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case int32:
		return int64(t)
	case float64:
		return int64(t)
	case float32:
		return int64(t)
	case string:
		var n int64
		fmt.Sscanf(t, "%d", &n)
		return n
	}
	return 0
}

// toAnySlice converts []map[string]any to []any for protobuf structpb compatibility.
// structpb.NewStruct() silently fails on []map[string]any but works with []any.
func toAnySlice(s []map[string]any) []any {
	result := make([]any, len(s))
	for i, m := range s {
		result[i] = m
	}
	return result
}

func toFloat64(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int64:
		return float64(t)
	case int:
		return float64(t)
	case string:
		var f float64
		fmt.Sscanf(t, "%f", &f)
		return f
	}
	return 0
}
