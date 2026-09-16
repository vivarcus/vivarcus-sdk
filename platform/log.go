package platform

// Log writes a guest log line via the host (vivarcus:platform/log.log).
// The reactor codegen wires this to the imported host function.
var LogFunc func(level, message string)

// LogInfo logs at info level.
func LogInfo(message string) {
	if LogFunc != nil {
		LogFunc("info", message)
	}
}

// LogWarn logs at warn level.
func LogWarn(message string) {
	if LogFunc != nil {
		LogFunc("warn", message)
	}
}

// LogError logs at error level.
func LogError(message string) {
	if LogFunc != nil {
		LogFunc("error", message)
	}
}
