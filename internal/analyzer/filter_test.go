package analyzer

import (
	"reflect"
	"strings"
	"testing"
)

func TestSlimStackTrace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "node stack drops system frames",
			input: `
ReferenceError: x is not defined
    at Object.<anonymous> (C:\proj\node_modules\pkg\index.js:1:1)
    at Module._compile (node:internal/modules/cjs/loader:1529:14)
    at Object.<anonymous> (C:\proj\src\index.js:10:5)
`,
			want: `
ReferenceError: x is not defined
    at Object.<anonymous> (C:\proj\src\index.js:10:5)
`,
		},
		{
			name: "python traceback drops runpy and lib/python3",
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
			want: `
Traceback (most recent call last):
  File "C:\proj\app\main.py", line 5, in <module>
    foo()
  File "C:\proj\app\utils.py", line 12, in foo
    raise ValueError("boom")
ValueError: boom
`,
		},
		{
			name: "drops async_hooks and v8 frames",
			input: `
Error: boom
    at hook (node:async_hooks:123:45)
    at native (v8::internal::Isolate:99:1)
    at C:\proj\src\app.js:7:3
`,
			want: `
Error: boom
    at C:\proj\src\app.js:7:3
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SlimStackTrace(lines(tt.input))
			want := lines(tt.want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("SlimStackTrace() =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
			}
		})
	}
}

func TestIsBlacklistedLine(t *testing.T) {
	blacklisted := []string{
		`    at Module._compile (node:internal/modules/cjs/loader:1529:14)`,
		`  File "C:\Python312\Lib\python3.12\runpy.py", line 88, in _run_code`,
		`    at hook (node:async_hooks:123:45)`,
		`    at native (v8::internal::Isolate:99:1)`,
	}
	for _, line := range blacklisted {
		if !isBlacklistedLine(line) {
			t.Errorf("expected blacklisted line: %q", line)
		}
	}
}
