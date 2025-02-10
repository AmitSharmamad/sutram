package shared

type LowerThan struct{}

func (o *LowerThan) Operate(lhs float64, rhs float64) bool {
	return lhs < rhs
}
