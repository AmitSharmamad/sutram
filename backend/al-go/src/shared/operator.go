package shared

type Operator interface {
	Operate(lhs float64, rhs float64) bool
}

type RaiseType string

// Up or Down
type RaiseBy struct {
	Value float64
	Type  RaiseType
}

const (
	PERCENTAGE RaiseType = "PERCENTAGE"
	ABSOLUTE   RaiseType = "ABSOLUTE"
)
