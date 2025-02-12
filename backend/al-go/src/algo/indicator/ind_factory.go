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
	factory.indicators["WMA"] = &WMA{}
	factory.indicators["DEMA"] = &DEMA{}
	factory.indicators["TEMA"] = &TEMA{}
	factory.indicators["HMA"] = &HMA{}
	factory.indicators["JMA"] = &JMA{}
	factory.indicators["MCGINLEY"] = &McGinley{}
	factory.indicators["VIDYA"] = &VIDYA{}
	factory.indicators["ALMA"] = &ALMA{}
	factory.indicators["ZLEMA"] = &ZLEMA{}
	factory.indicators["RMA"] = &RMA{}
	factory.indicators["AWMA"] = &AnchoredWMA{}

	// Advanced Moving Averages
	factory.indicators["HULLVOL"] = &HullVolume{}
	factory.indicators["HMACHAIKIN"] = &HMAChaikin{}
	factory.indicators["VWHB"] = &VWHB{}

	// Momentum Indicators
	factory.indicators["RSI"] = &RSI{}
	factory.indicators["MACD"] = &MACD{}
	factory.indicators["CCI"] = &CCI{}
	factory.indicators["CMO"] = &CMO{}
	factory.indicators["MOM"] = &MOM{}
	factory.indicators["ROC"] = &ROC{}
	factory.indicators["TRIX"] = &TRIX{}
	factory.indicators["TSI"] = &TSI{}
	factory.indicators["UO"] = &UO{}
	factory.indicators["QQE"] = &QQE{}
	factory.indicators["DV2"] = &DV2{}
	factory.indicators["DMI"] = &DMI{}
	factory.indicators["RAVI"] = &RAVI{}
	factory.indicators["PM"] = &PM{}
	factory.indicators["RCI"] = &RCI{}
	factory.indicators["FORCEINDEX"] = &ForceIndex{}
	factory.indicators["AWM"] = &AWM{}

	// Trend Indicators
	factory.indicators["ADX"] = &ADX{}
	factory.indicators["AO"] = &AO{}
	factory.indicators["CHOP"] = &CHOP{}
	factory.indicators["COG"] = &COG{}
	factory.indicators["DPO"] = &DPO{}
	factory.indicators["EMV"] = &EMV{}
	factory.indicators["ELDERRAY"] = &ElderRay{}
	factory.indicators["KAMA"] = &KAMA{}
	factory.indicators["KST"] = &KST{}
	factory.indicators["PPO"] = &PPO{}
	factory.indicators["SUPERTREND"] = &SuperTrend{}
	factory.indicators["VORTEX"] = &Vortex{}
	factory.indicators["TRENDFLEX"] = &TrendFlex{}
	// factory.indicators["APZ"] = &APZ{}
	factory.indicators["STC"] = &STC{}
	factory.indicators["WAVETREND"] = &WaveTrend{}

	// Volatility Indicators
	factory.indicators["ATR"] = &ATR{}
	factory.indicators["BB"] = &BollingerBands{}
	factory.indicators["MASSINDEX"] = &MassIndex{}

	// Volume Indicators
	factory.indicators["VWAP"] = &VWAP{}
	factory.indicators["VPT"] = &VPT{}

	// Statistical/Cyclical Indicators
	factory.indicators["ENTROPY"] = &Entropy{}
	factory.indicators["FISHER"] = &Fisher{}
	factory.indicators["SKEW"] = &Skew{}
	factory.indicators["STOCH"] = &Stochastic{}
	factory.indicators["WILLIAMSR"] = &WilliamsR{}
	factory.indicators["MESASINE"] = &MesaSineWave{}

	// Price Transform
	factory.indicators["HA"] = &HA{}

	// Utility Indicator
	factory.indicators["Number"] = &IndNumber{}

	return factory
}

var indFactory = NewFactory()

func GetIndFactory(key string) IndicatorCalculator {
	return indFactory.indicators[key]
}
