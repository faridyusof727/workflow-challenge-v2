package condition

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ExprReplacePlaceholderByMap replaces placeholders in an expression string with values from a struct.
// It iterates through the fields of the input struct, using JSON tags as placeholders.
// For each field with a JSON tag, it replaces the corresponding placeholder in the expression.
// Special handling is provided for the "operator" field to convert it to an expression.
// Returns the modified expression string with all placeholders replaced by their corresponding values.
func ExprReplacePlaceholderByMap(inputs Inputs) string {
	exprString := inputs.Expression

	t := reflect.TypeOf(inputs)
	v := reflect.ValueOf(inputs)

	for i := range t.NumField() {

		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		if jsonTag == "" || jsonTag == "expression" {
			continue // skip fields without a json tag
		}

		placeholder := "{{" + jsonTag + "}}"

		// Get the field value as an interface{}
		fieldValue := v.Field(i).Interface()

		if jsonTag == "operator" {
			// Type check before type assertion
			if op, ok := fieldValue.(Operator); ok {
				exprString = ExprReplacePlaceholder(exprString, placeholder, op.ToExpr())
			} else {
				continue
			}
		} else {
			exprString = ExprReplacePlaceholder(exprString, placeholder, fieldValue)
		}

	}

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
