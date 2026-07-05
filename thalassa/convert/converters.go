package convert

func Ptr[T any](v T) *T {
	return &v
}

func StringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func ConvertToStringSlice(v any) []string {
	if v == nil {
		return []string{}
	}

	result := []string{}
	for _, item := range v.([]any) {
		result = append(result, item.(string))
	}

	return result
}

func ConvertToMap(v any) map[string]string {
	if v == nil {
		return map[string]string{}
	}

	result := map[string]string{}
	for key, value := range v.(map[string]any) {
		result[key] = value.(string)
	}

	return result
}

func ConvertToInt32Slice(v any) []int32 {
	if v == nil {
		return []int32{}
	}

	result := []int32{}
	for _, item := range v.([]any) {
		result = append(result, int32(item.(int)))
	}

	return result
}
