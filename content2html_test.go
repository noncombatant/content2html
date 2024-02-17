package content2html

import "testing"

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
