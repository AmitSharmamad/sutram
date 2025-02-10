package shared

type DownBy RaiseBy

func (o *DownBy) Operate(lhs float64, rhs float64) bool {
	if o.Type == PERCENTAGE {
		return lhs < rhs*(1-o.Value/100)
	} else if o.Type == ABSOLUTE {
		return lhs < rhs-o.Value
	}
	return false
}
