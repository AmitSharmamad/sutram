package shared

type CrossesBelow struct{}

func (*CrossesBelow) Operate(left float64, right float64) bool {
	return left < right
}
