package parser

type IndicatorFactory struct {
	indicators map[string]IndicatorCalculator
}

func NewFactory() *IndicatorFactory {
	factory := &IndicatorFactory{
		indicators: make(map[string]IndicatorCalculator),
	}

	// Price indicators
	factory.indicators["OPEN"] = &Open{}
	factory.indicators["HIGH"] = &High{}
	factory.indicators["LOW"] = &Low{}
	factory.indicators["CLOSE"] = &Close{}
	factory.indicators["VOLUME"] = &Volume{}

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
