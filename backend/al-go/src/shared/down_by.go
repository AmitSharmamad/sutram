package shared

type DownBy struct{}

func (o *DownBy) Operate(lhs float64, rhs float64) bool {
	return lhs < rhs
}
