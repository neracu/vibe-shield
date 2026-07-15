package analyzer

import "regexp"

var (
	rePyTraceback = regexp.MustCompile(`^Traceback \(most recent call last\):`)
	rePyFile      = regexp.MustCompile(`^\s*File "(.+?)", line (\d+),`)
	rePyError     = regexp.MustCompile(`^(\w+(?:Error|Exception|Interrupt)): (.+)$`)

	reNodeError      = regexp.MustCompile(`(?:^|:\s*)((?:\w+)?(?:Error|Exception)): (.+)$`)
	reNodeStackParen = regexp.MustCompile(`at\s+(?:.*?\s+)?\((.+):(\d+):\d+\)`)
	reNodeStackBare  = regexp.MustCompile(`at\s+([^(\s]+):(\d+):\d+`)
)
