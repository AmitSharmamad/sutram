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

	// Technical Indicators
	factory.indicators["SUPERTREND"] = &Supertrend{}
	factory.indicators["MACD"] = &MACD{}
	factory.indicators["VTX"] = &Vortex{}
	factory.indicators["BOLLINGER"] = &BollingerBands{}
	factory.indicators["LBB"] = &LowerBollingerBands{}
	factory.indicators["MBB"] = &MiddleBollingerBands{}
	factory.indicators["UBB"] = &UpperBollingerBands{}
	factory.indicators["DMI"] = &DMI{}
	factory.indicators["KELTNER"] = &KeltnerChannels{}
	factory.indicators["PSAR"] = &ParabolicSAR{}
	factory.indicators["FISHER"] = &FisherTransform{}
	factory.indicators["MFI"] = &MFI{}
	factory.indicators["VWAP"] = &VWAP{}
	factory.indicators["TSI"] = &TSI{}
	factory.indicators["VO"] = &VolumeOscillator{}
	factory.indicators["UO"] = &UltimateOscillator{}

	// Utility Indicator
	factory.indicators["Number"] = &IndNumber{}

	return factory
}

var indFactory = NewFactory()

func GetIndFactory(key string) IndicatorCalculator {
	return indFactory.indicators[key]
}
