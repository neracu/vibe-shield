package prompt

const maxTailLogs = 10

// TailLogs returns the last up-to-10 lines from captured stdout.
func TailLogs(lines []string) []string {
	if len(lines) == 0 {
		return nil
	}
	start := len(lines) - maxTailLogs
	if start < 0 {
		start = 0
	}
	return lines[start:]
}
