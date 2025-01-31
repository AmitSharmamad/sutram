package shared

type Operator interface {
	Operate(lhs float64, rhs float64) bool
}
