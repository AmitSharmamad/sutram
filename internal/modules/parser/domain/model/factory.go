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
	factory.indicators["ADX"] = &Adx{}
	factory.indicators["ATR"] = &Atr{}
	factory.indicators["AVGPRICE"] = &AvgPrice{}
	factory.indicators["CCI"] = &Cci{}
	factory.indicators["DEMA"] = &Dema{}
	factory.indicators["MAMA"] = &Mama{}
	factory.indicators["MINUSDI"] = &MinusDI{}
	factory.indicators["MOM"] = &Mom{}
	factory.indicators["PLUSDI"] = &PlusDI{}
	factory.indicators["ROC"] = &Roc{}
	factory.indicators["RSI"] = &Rsi{}
	factory.indicators["TEMA"] = &Tema{}
	factory.indicators["WILLR"] = &WillR{}
	factory.indicators["WMA"] = &Wma{}

	// Utility Indicator
	factory.indicators["Number"] = &IndNumber{}

	return factory
}

var indFactory = NewFactory()

func GetIndFactory(key string) IndicatorCalculator {
	return indFactory.indicators[key]
}
