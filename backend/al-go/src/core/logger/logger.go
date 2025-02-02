package logger

func Log(values ...interface{}) {
	for _, v := range values {
		println(v)
	}
}
