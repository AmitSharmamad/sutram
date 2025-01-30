package indicator

type IndicatorFactory struct {
	indicators map[string]IndicatorCalculator
}

func NewFactory() *IndicatorFactory {
	factory := &IndicatorFactory{
		indicators: make(map[string]IndicatorCalculator),
	}

	factory.indicators["SMA"] = &SMA{}
	factory.indicators["EMA"] = &EMA{}
	factory.indicators["Number"] = &IndNumber{}

	return factory
}

var factory = NewFactory()

func Get(key string) IndicatorCalculator {
	return factory.indicators[key]
}
