package condition

type Operator string

const (
	GreaterThanOperator Operator = "greater_than"
	LessThanOperator    Operator = "less_than"
	EqualToOperator     Operator = "equals"
	IsAtLeastOperator   Operator = "greater_than_or_equal"
	IsAtMostOperator    Operator = "less_than_or_equal"
)

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
