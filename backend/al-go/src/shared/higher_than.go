package shared

type HigherThan struct{}

func (o *HigherThan) Operate(lhs float64, rhs float64) bool {
	return lhs > rhs
}
