package operator

func Init() {
	println("Initialising Operators")

	InitFactory()

	crossesAbove := &CrossesAbove{}
	crossesBelow := &CrossesBelow{}

	Register("CrossesAbove", crossesAbove)
	Register("CrossesBelow", crossesBelow)
}
