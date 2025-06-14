package condition

import (
	"fmt"
	"strconv"
	"strings"
)

// ExprReplacePlaceholderByMap replaces placeholders in an expression string with values from a map.
// It iterates through the map, replacing placeholders of the format {{key}} with their corresponding values.
// Additionally, it replaces the {{operator}} placeholder with the string representation of the given operator.
// Returns the modified expression string with all placeholders replaced.
func ExprReplacePlaceholderByMap(exprString string, maps map[string]interface{}, operator Operator) string {
	for key, value := range maps {
		placeholder := "{{" + key + "}}"
		exprString = ExprReplacePlaceholder(exprString, placeholder, value)
	}

	exprString = ExprReplacePlaceholder(exprString, "{{operator}}", operator.ToExpr())

	return exprString
}

// ExprReplacePlaceholder replaces a specific placeholder in an expression string with a formatted value.
// It uses the formatValue function to convert the value to a string representation
// and then replaces all occurrences of the placeholder with that formatted value.
// Returns the modified expression string with the placeholder replaced.
func ExprReplacePlaceholder(expr, placeholder string, value interface{}) string {
	return strings.ReplaceAll(expr, placeholder, formatValue(value))
}

// formatValue converts an interface{} value to its string representation.
// It handles various types including string, float64, int, bool, float32, and int64.
// For unsupported types, it uses fmt.Sprintf as a fallback to convert the value to a string.
// Returns the string representation of the input value.
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		return fmt.Sprintf("%v", v) // Fallback for other types
	}
}
