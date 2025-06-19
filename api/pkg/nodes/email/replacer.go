package email

import (
	"fmt"
	"strings"
)

// TmplReplacePlaceholderByMap replaces placeholders in the format {{key}} within a template string.
// tmpl is the template string containing placeholders to be replaced.
// args is a map containing key-value pairs where keys match placeholder names.
// Returns the template with all matching placeholders replaced by their corresponding values.
func TmplReplacePlaceholderByMap(tmpl string, args map[string]any) string {
	for key, value := range args {
		placeholder := fmt.Sprintf("{{%s}}", key)
		replacement := fmt.Sprintf("%v", value)
		tmpl = strings.ReplaceAll(tmpl, placeholder, replacement)
	}

	return tmpl
}