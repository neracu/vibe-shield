package prompt

const maxTailLogs = 10
const maxTailStderr = 20

// TailLogs returns the last up-to-10 lines from captured stdout.
func TailLogs(lines []string) []string {
	return tailLines(lines, maxTailLogs)
}

// TailStderr returns the last up-to-20 lines from captured stderr.
func TailStderr(lines []string) []string {
	return tailLines(lines, maxTailStderr)
}

func tailLines(lines []string, max int) []string {
	if len(lines) == 0 {
		return nil
	}
	start := len(lines) - max
	if start < 0 {
		start = 0
	}
	return lines[start:]
}
