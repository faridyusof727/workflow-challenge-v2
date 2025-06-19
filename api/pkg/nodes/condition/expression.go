package condition

import (
	"fmt"
	"strings"
)

// ExprReplacePlaceholderByMap replaces placeholders in the format {{key}} within an expression string.
// conditionExpression is the template string containing placeholders to be replaced.
// args is a map containing key-value pairs where keys match placeholder names.
// Special handling is applied for OperatorKey values which are converted using ToExpr().
// Returns the expression with all matching placeholders replaced by their corresponding values.
func ExprReplacePlaceholderByMap(conditionExpression string, args map[string]any) string {
	result := conditionExpression

	for key, value := range args {
		placeholder := fmt.Sprintf("{{%s}}", key)

		var replacement string
		if key == OperatorKey {
			op := Operator(value.(string))
			replacement = fmt.Sprintf("%v", op.ToExpr())
		} else {
			replacement = fmt.Sprintf("%v", value)
		}
		result = strings.ReplaceAll(result, placeholder, replacement)
	}

	return result
}
