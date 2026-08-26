package cli

const (
	dimStyle       = "\x1b[2m"
	highlightStyle = "\x1b[48;5;238m\x1b[38;5;255m"
	resetStyle     = "\x1b[0m"
)

func dim(text string) string {
	return dimStyle + text + resetStyle
}
func highlight(command string) string {
	return highlightStyle + " " + command + " " + resetStyle
}
