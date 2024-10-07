package content2html

import (
	"bytes"
	"testing"
)

func TestGetHTMLPathname(t *testing.T) {
	tests := []struct {
		name          string
		pathname      string
		outputDirname string
		want          string
		wantError     bool
	}{
		{
			name:          "content sibling",
			pathname:      "foo.content",
			outputDirname: "",
			want:          "foo.html",
		},
		{
			name:          "content with output dir",
			pathname:      "foo.content",
			outputDirname: "out",
			want:          "out/foo.html",
		},
		{
			name:          "nested with output dir",
			pathname:      "foo/bar/quux.content",
			outputDirname: "out",
			want:          "out/foo/bar/quux.html",
		},
		{
			name:          "sibling collision",
			pathname:      "foo/bar/quux.html",
			outputDirname: "",
			want:          "",
			wantError:     true,
		},
	}
	for _, tt := range tests {
		result, e := GetHTMLPathname(tt.pathname, tt.outputDirname)
		if result != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, result, tt.want)
		}
		if tt.wantError && e == nil {
			t.Errorf("%s: did not get error", tt.name)
		}
	}
}

func TestGetHairSpaces(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "one line with spaces",
			input: "hello, this is — a test",
			want:  "hello, this is\u200A—\u200Aa test",
		},
		{
			name:  "one line with left space only",
			input: "hello, this is —a test",
			want:  "hello, this is\u200A—\u200Aa test",
		},
		{
			name:  "one line with right space only",
			input: "hello, this is— a test",
			want:  "hello, this is\u200A—\u200Aa test",
		},
		{
			name:  "beginning of line, space",
			input: "— hello, this is a test",
			want:  "\u200A—\u200ahello, this is a test",
		},
		{
			name:  "beginning of line, no space",
			input: "—hello, this is a test",
			want:  "\u200A—\u200ahello, this is a test",
		},
		{
			name:  "end of line, space",
			input: "hello, this is a test —",
			want:  "hello, this is a test\u200A—\u200a",
		},
		{
			name:  "end of line, no space",
			input: "hello, this is a test—",
			want:  "hello, this is a test\u200A—\u200a",
		},
		{
			name: "multiline 1",
			input: `hello, this is a test—
			of the emergency goatcast system`,
			want: "hello, this is a test\u200A—\u200aof the emergency goatcast system",
		},
		{
			name: "multiline 2",
			input: `hello, this is a test —
			of the emergency goatcast system`,
			want: "hello, this is a test\u200A—\u200aof the emergency goatcast system",
		},
		{
			name: "multiline 3",
			input: `— hello, this is a test —
			of the emergency goatcast system`,
			want: "\u200A—\u200ahello, this is a test\u200A—\u200aof the emergency goatcast system",
		},
		{
			name: "multiline 4",
			input: `—hello, this is a test—
			of the emergency goatcast system`,
			want: "\u200A—\u200ahello, this is a test\u200A—\u200aof the emergency goatcast system",
		},
	}
	for _, tt := range tests {
		result := useHairSpaces([]byte(tt.input))
		if !bytes.Equal(result, []byte(tt.want)) {
			t.Errorf("%s: got %q, want %q", tt.name, result, tt.want)
		}
	}
}
