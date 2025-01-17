package indicator

func Init() {
	println("Initialising Indicators")

	InitFactory()

	sma := &SMA{}
	ema := &EMA{}
	number := &IndNumber{}

	Register("SMA", sma)
	Register("EMA", ema)
	Register("Number", number)

}
