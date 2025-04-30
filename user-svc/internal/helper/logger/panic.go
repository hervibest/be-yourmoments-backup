package logger

func (l *logImpl) Panic(message interface{}, options ...Options) {
	l.stderr.Panicln(message)
}

func (l *logImpl) CustomPanic(title string, message interface{}, options ...Options) {
	msg := []interface{}{}
	msg = append(msg, title)
	msg = append(msg, message)

	l.stderr.Panicln(msg...)
}
