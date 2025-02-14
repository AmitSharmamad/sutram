package shared

type EqualTo struct{}

func (o *EqualTo) Operate(lhs float64, rhs float64) bool {
	return lhs == rhs
}
