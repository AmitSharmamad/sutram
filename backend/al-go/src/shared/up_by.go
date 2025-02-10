package shared

type UpBy struct{}

func (o *UpBy) Operate(lhs float64, rhs float64) bool {
	return lhs > rhs
}
