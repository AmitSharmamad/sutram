package indicator

type IndicatorFactory struct {
	indicators map[string]IndicatorCalculator
}

func NewFactory() *IndicatorFactory {
	factory := &IndicatorFactory{
		indicators: make(map[string]IndicatorCalculator),
	}

	// Basic Moving Averages
	factory.indicators["SMA"] = &SMA{}
	factory.indicators["EMA"] = &EMA{}

	// Utility Indicator
	factory.indicators["Number"] = &IndNumber{}

	return factory
}

var indFactory = NewFactory()

func GetIndFactory(key string) IndicatorCalculator {
	return indFactory.indicators[key]
}
