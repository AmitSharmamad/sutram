package indicator

var funcs map[string]*IndicatorCalculator

func InitFactory() {
	funcs = make(map[string]*IndicatorCalculator, 0)
}

func Register(key string, calculator IndicatorCalculator) {
	funcs[key] = &calculator
}

func Get(key string) *IndicatorCalculator {
	println("Getting indicator ", key)
	return funcs[key]
}
