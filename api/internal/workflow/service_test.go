package workflow

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	data := []byte(`{
  "formData": {
    "name": "asd",
    "email": "asd@asd.casd",
    "city": "Sydney",
    "operator": "greater_than",
    "threshold": 25
  },
  "condition": {
    "operator": "greater_than",
    "threshold": 25
  }
}`)

	var jsonData map[string]any
	json.Unmarshal(data, &jsonData)

	flattened := make(map[string]any)
	flattenMap(jsonData, "", flattened)

	for key, value := range flattened {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func flattenMap(data map[string]any, prefix string, result map[string]any) {
	for key, value := range data {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		if nested, ok := value.(map[string]any); ok {
			flattenMap(nested, newKey, result)
		} else {
			result[newKey] = value
		}
	}
}
