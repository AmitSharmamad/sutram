package operator

type OperatorFactory struct {
	operators map[string]Operator
}

func NewFactory() *OperatorFactory {
	factory := &OperatorFactory{
		operators: make(map[string]Operator),
	}

	factory.operators["CrossesAbove"] = &CrossesAbove{}
	factory.operators["CrossesBelow"] = &CrossesBelow{}

	return factory
}

var factory = NewFactory()

func Get(key string) Operator {
	return factory.operators[key]
}
