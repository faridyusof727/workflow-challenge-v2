package condition

import "fmt"

type Operator string

const (
	GreaterThanOperator Operator = "greater_than"
	LessThanOperator    Operator = "less_than"
	EqualToOperator     Operator = "equals"
	IsAtLeastOperator   Operator = "greater_than_or_equal"
	IsAtMostOperator    Operator = "less_than_or_equal"
)

// Validate checks if the Operator is a valid predefined operator.
// Returns nil if the operator is valid, otherwise returns an error with details about the invalid operator.
func (o Operator) Validate() error {
	switch o {
	case GreaterThanOperator, LessThanOperator, EqualToOperator, IsAtLeastOperator, IsAtMostOperator:
		return nil
	default:
		return fmt.Errorf("invalid operator: %s", o)
	}
}

// ToExpr converts the Operator to its corresponding symbolic representation.
// Returns the string symbol for the operator (e.g., ">" for GreaterThanOperator).
// Returns an empty string if the operator is not recognized.
func (o Operator) ToExpr() string {
	switch o {
	case GreaterThanOperator:
		return ">"
	case LessThanOperator:
		return "<"
	case EqualToOperator:
		return "=="
	case IsAtLeastOperator:
		return ">="
	case IsAtMostOperator:
		return "<="
	default:
		return ""
	}
}
