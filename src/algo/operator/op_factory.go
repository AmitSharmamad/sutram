package operator

type OperatorFactory struct {
	operators map[string]Operator
}

func NewOperatorFactory() *OperatorFactory {
	factory := &OperatorFactory{
		operators: make(map[string]Operator),
	}

	factory.operators["CrossesAbove"] = &CrossesAbove{}
	factory.operators["CrossesBelow"] = &CrossesBelow{}
	factory.operators["UpBy"] = &UpBy{}
	factory.operators["DownBy"] = &DownBy{}
	factory.operators["EqualTo"] = &EqualTo{}
	factory.operators["LowerThan"] = &LowerThan{}
	factory.operators["HigherThan"] = &HigherThan{}
	factory.operators["LowerThanEqualTo"] = &LowerThanEqualTo{}
	factory.operators["HigherThanEqualTo"] = &HigherThanEqualTo{}

	return factory
}

var opFactory = NewOperatorFactory()

func GetOperator(key string) Operator {
	return opFactory.operators[key]
}
