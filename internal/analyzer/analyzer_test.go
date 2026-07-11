package analyzer

import (
	"reflect"
	"strings"
	"testing"
)

func lines(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func TestAnalyzeStderr(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOK  bool
		wantErr *DetectedError
	}{
		{
			name: "node reference error skips node_modules",
			input: `
ReferenceError: x is not defined
    at Object.<anonymous> (C:\proj\node_modules\pkg\index.js:1:1)
    at Module._compile (node:internal/modules/cjs/loader:1529:14)
    at Object.<anonymous> (C:\proj\src\index.js:10:5)
`,
			wantOK: true,
			wantErr: &DetectedError{
				ErrorType:    "ReferenceError",
				ErrorMessage: "x is not defined",
				FilePath:     `C:\proj\src\index.js`,
				LineNumber:   10,
				StackTrace: []string{
					"ReferenceError: x is not defined",
					`    at Object.<anonymous> (C:\proj\src\index.js:10:5)`,
				},
			},
		},
		{
			name: "node type error bare stack frame",
			input: `
TypeError: Cannot read properties of undefined
    at C:\proj\path\file.ts:42:5
    at processTicksAndRejections (node:internal/process/task_queues:95:5)
`,
			wantOK: true,
			wantErr: &DetectedError{
				ErrorType:    "TypeError",
				ErrorMessage: "Cannot read properties of undefined",
				FilePath:     `C:\proj\path\file.ts`,
				LineNumber:   42,
				StackTrace: []string{
					"TypeError: Cannot read properties of undefined",
					`    at C:\proj\path\file.ts:42:5`,
				},
			},
		},
		{
			name: "python traceback skips lib/python",
			input: `
Traceback (most recent call last):
  File "C:\proj\app\main.py", line 5, in <module>
    foo()
  File "C:\Python312\Lib\python3.12\runpy.py", line 88, in _run_code
    exec(code, run_globals)
  File "C:\proj\app\utils.py", line 12, in foo
    raise ValueError("boom")
ValueError: boom
`,
			wantOK: true,
			wantErr: &DetectedError{
				ErrorType:    "ValueError",
				ErrorMessage: "boom",
				FilePath:     `C:\proj\app\utils.py`,
				LineNumber:   12,
				StackTrace: []string{
					"Traceback (most recent call last):",
					`  File "C:\proj\app\main.py", line 5, in <module>`,
					"    foo()",
					`  File "C:\proj\app\utils.py", line 12, in foo`,
					`    raise ValueError("boom")`,
					"ValueError: boom",
				},
			},
		},
		{
			name: "only system paths returns false",
			input: `
ReferenceError: x is not defined
    at Object.<anonymous> (C:\proj\node_modules\pkg\index.js:1:1)
    at Module._compile (node:internal/modules/cjs/loader:1529:14)
`,
			wantOK: false,
		},
		{
			name: "unrelated stderr noise",
			input: `
[INFO] server starting on port 3000
[WARN] deprecated API usage
`,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := AnalyzeStderr(lines(tt.input))
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}
			if got.ErrorType != tt.wantErr.ErrorType {
				t.Errorf("ErrorType = %q, want %q", got.ErrorType, tt.wantErr.ErrorType)
			}
			if got.ErrorMessage != tt.wantErr.ErrorMessage {
				t.Errorf("ErrorMessage = %q, want %q", got.ErrorMessage, tt.wantErr.ErrorMessage)
			}
			if got.FilePath != tt.wantErr.FilePath {
				t.Errorf("FilePath = %q, want %q", got.FilePath, tt.wantErr.FilePath)
			}
			if got.LineNumber != tt.wantErr.LineNumber {
				t.Errorf("LineNumber = %d, want %d", got.LineNumber, tt.wantErr.LineNumber)
			}
			if !reflect.DeepEqual(got.StackTrace, tt.wantErr.StackTrace) {
				t.Errorf("StackTrace =\n%s\nwant\n%s", strings.Join(got.StackTrace, "\n"), strings.Join(tt.wantErr.StackTrace, "\n"))
			}
		})
	}
}

func TestIsSystemPath(t *testing.T) {
	system := []string{
		`/app/node_modules/foo/index.js`,
		`node:internal/modules/cjs/loader`,
		`webpack-internal:///./src/index.js`,
		`/usr/lib/python3.12/os.py`,
		`/venv/lib/python3.12/site-packages/requests/__init__.py`,
		`/usr/lib/python3/dist-packages/pkg.py`,
		`C:\Python312\Lib\python3.12\runpy.py`,
		`node:async_hooks`,
		`v8::internal::Isolate`,
	}
	for _, p := range system {
		if !isSystemPath(p) {
			t.Errorf("expected system path: %q", p)
		}
	}

	user := []string{
		`C:\proj\src\index.js`,
		`/home/user/app/main.py`,
	}
	for _, p := range user {
		if isSystemPath(p) {
			t.Errorf("expected user path: %q", p)
		}
	}
}
