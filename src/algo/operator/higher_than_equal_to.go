package operator

type HigherThanEqualTo struct{}

func (o *HigherThanEqualTo) Operate(lhs float64, rhs float64) bool {
	return lhs >= rhs
}
