package thalassa

func Ptr[T any](v T) *T {
	return &v
}

func convertToStringSlice(v interface{}) []string {
	if v == nil {
		return []string{}
	}

	result := []string{}
	for _, item := range v.([]interface{}) {
		result = append(result, item.(string))
	}

	return result
}

func convertToMap(v interface{}) map[string]string {
	if v == nil {
		return map[string]string{}
	}

	result := map[string]string{}
	for key, value := range v.(map[string]interface{}) {
		result[key] = value.(string)
	}

	return result
}
