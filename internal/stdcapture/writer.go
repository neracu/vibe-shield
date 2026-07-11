package stdcapture

import (
	"bytes"
	"os"
	"strings"
)

// Writer intercepts stderr: streams each line to os.Stderr in real time
// and accumulates lines for later analysis.
type Writer struct {
	lines []string
	buf   []byte
}

func New() *Writer {
	return &Writer{}
}

func (w *Writer) Write(p []byte) (int, error) {
	n := len(p)
	w.buf = append(w.buf, p...)

	for {
		idx := bytes.IndexByte(w.buf, '\n')
		if idx < 0 {
			break
		}
		line := w.buf[:idx+1]
		w.buf = w.buf[idx+1:]

		if _, err := os.Stderr.Write(line); err != nil {
			return n, err
		}
		w.lines = append(w.lines, strings.TrimSuffix(string(line), "\n"))
	}

	return n, nil
}

// Flush emits and stores any trailing bytes not terminated by a newline.
func (w *Writer) Flush() error {
	if len(w.buf) == 0 {
		return nil
	}
	if _, err := os.Stderr.Write(w.buf); err != nil {
		return err
	}
	w.lines = append(w.lines, string(w.buf))
	w.buf = w.buf[:0]
	return nil
}

func (w *Writer) Lines() []string {
	return w.lines
}
