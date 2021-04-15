package web

// nilLogger represents a logger that does not print anything
type nilLogger struct {
}

// Printf prints nothing
func (logger *nilLogger) Printf(string, ...interface{}) {
}
