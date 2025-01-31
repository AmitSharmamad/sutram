package shared

type OperatorFactory struct {
	operators map[string]Operator
}

func NewOperatorFactory() *OperatorFactory {
	factory := &OperatorFactory{
		operators: make(map[string]Operator),
	}

	factory.operators["CrossesAbove"] = &CrossesAbove{}
	factory.operators["CrossesBelow"] = &CrossesBelow{}

	return factory
}

var opFactory = NewOperatorFactory()

func GetOperator(key string) Operator {
	return opFactory.operators[key]
}
