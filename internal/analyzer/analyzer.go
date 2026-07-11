package analyzer

import (
	"strconv"
	"strings"
)

type DetectedError struct {
	ErrorType    string
	ErrorMessage string
	FilePath     string
	LineNumber   int
	StackTrace   []string
}

func AnalyzeStderr(buffer []string) (*DetectedError, bool) {
	for i, line := range buffer {
		if rePyTraceback.MatchString(line) {
			if detected, ok := parsePythonTraceback(buffer[i:]); ok {
				return detected, true
			}
		}
	}

	for i, line := range buffer {
		if m := reNodeError.FindStringSubmatch(line); m != nil {
			if detected, ok := parseNodeError(line, m[1], m[2], buffer[i+1:]); ok {
				return detected, true
			}
		}
	}

	return nil, false
}

func parsePythonTraceback(lines []string) (*DetectedError, bool) {
	traceBlock := []string{lines[0]}

	var (
		errType    string
		errMessage string
		userFile   string
		userLine   int
	)

	for _, line := range lines[1:] {
		if line == "" {
			break
		}
		traceBlock = append(traceBlock, line)

		if m := rePyFile.FindStringSubmatch(line); m != nil {
			if !isSystemPath(m[1]) {
				userFile = m[1]
				userLine, _ = strconv.Atoi(m[2])
			}
			continue
		}
		if m := rePyError.FindStringSubmatch(line); m != nil {
			errType = m[1]
			errMessage = m[2]
			break
		}
		if !strings.HasPrefix(line, "  ") {
			break
		}
	}

	if errType == "" || userFile == "" {
		return nil, false
	}

	return &DetectedError{
		ErrorType:    errType,
		ErrorMessage: errMessage,
		FilePath:     userFile,
		LineNumber:   userLine,
		StackTrace:   SlimStackTrace(traceBlock),
	}, true
}

func parseNodeError(errorLine, errType, errMessage string, following []string) (*DetectedError, bool) {
	traceBlock := []string{errorLine}

	var (
		userFile string
		userLine int
	)

	for _, line := range following {
		if reNodeError.MatchString(line) {
			break
		}
		path, lineNum, ok := extractNodeFrame(line)
		if !ok {
			break
		}
		traceBlock = append(traceBlock, line)
		if !isSystemPath(path) {
			userFile = path
			userLine = lineNum
		}
	}

	if userFile == "" {
		return nil, false
	}

	return &DetectedError{
		ErrorType:    errType,
		ErrorMessage: errMessage,
		FilePath:     userFile,
		LineNumber:   userLine,
		StackTrace:   SlimStackTrace(traceBlock),
	}, true
}

func extractNodeFrame(line string) (string, int, bool) {
	if m := reNodeStackParen.FindStringSubmatch(line); m != nil {
		n, _ := strconv.Atoi(m[2])
		return m[1], n, true
	}
	if m := reNodeStackBare.FindStringSubmatch(line); m != nil {
		n, _ := strconv.Atoi(m[2])
		return m[1], n, true
	}
	return "", 0, false
}
