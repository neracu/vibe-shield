package analyzer

import "strings"

var systemPathBlacklist = []string{
	"node_modules",
	"node:internal",
	"async_hooks",
	"v8",
	"lib/python3.",
	"site-packages",
	"runpy.py",
	"lib/python",
	"dist-packages",
	"webpack-internal",
}

func isBlacklistedLine(line string) bool {
	lower := strings.ToLower(line)
	for _, marker := range systemPathBlacklist {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func isSystemPath(path string) bool {
	return isBlacklistedLine(path)
}

func isStackFrameLine(line string) bool {
	if rePyFile.MatchString(line) {
		return true
	}
	_, _, ok := extractNodeFrame(line)
	return ok
}

func SlimStackTrace(lines []string) []string {
	out := make([]string, 0, len(lines))
	skipIndented := false

	for _, line := range lines {
		if isBlacklistedLine(line) {
			if isStackFrameLine(line) {
				skipIndented = true
			}
			continue
		}

		if isStackFrameLine(line) {
			if m := rePyFile.FindStringSubmatch(line); m != nil {
				if isSystemPath(m[1]) {
					skipIndented = true
					continue
				}
			}
			if path, _, ok := extractNodeFrame(line); ok && isSystemPath(path) {
				skipIndented = true
				continue
			}
			skipIndented = false
			out = append(out, line)
			continue
		}

		if skipIndented && strings.HasPrefix(line, "  ") {
			continue
		}

		skipIndented = false
		out = append(out, line)
	}

	return out
}
