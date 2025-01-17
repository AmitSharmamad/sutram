package operator

var funcs map[string]Operator

func InitFactory() {
	funcs = make(map[string]Operator, 0)
}

func Register(key string, operator Operator) {
	funcs[key] = operator
}

func Get(key string) Operator {
	return funcs[key]
}
