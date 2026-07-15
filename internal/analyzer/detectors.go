package analyzer

type Detector interface {
	Detect(buffer []string) (*DetectedError, bool)
}

var stderrDetectors = []Detector{
	pythonDetector{},
	nodeDetector{},
}

type pythonDetector struct{}

func (pythonDetector) Detect(buffer []string) (*DetectedError, bool) {
	for i, line := range buffer {
		if rePyTraceback.MatchString(line) {
			if detected, ok := parsePythonTraceback(buffer[i:]); ok {
				return detected, true
			}
		}
	}
	return nil, false
}

type nodeDetector struct{}

func (nodeDetector) Detect(buffer []string) (*DetectedError, bool) {
	for i, line := range buffer {
		if errType, errMessage, ok := matchNodeErrorLineWithContext(line, buffer[i+1:]); ok {
			if detected, ok := parseNodeError(line, errType, errMessage, buffer[i+1:]); ok {
				return detected, true
			}
		}
	}
	return nil, false
}
