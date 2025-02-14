package shared

type LowerThanEqualTo struct{}

func (o *LowerThanEqualTo) Operate(lhs float64, rhs float64) bool {
	return lhs <= rhs
}
