package operator

type CrossesAbove struct{}

func (*CrossesAbove) Operate(left float64, right float64) bool {
	return left > right
}
