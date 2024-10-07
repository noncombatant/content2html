// Package content2html implements a very simple HTML templating system for
// documents. It applies an html/template and fills it with the body text from
// an input file, discovering the document's <title> from its 1st <h1>. This
// way, you can write plain HTML, and generate complete documents with
// templatized boilerplate.
package content2html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

var minifier *minify.M

func init() {
	minifier = minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("image/svg+xml", svg.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
}

// Document represents a document with a raw HTML body, and a <title>.
// GenerateHTML will discover the Title.
//
// TODO: discover all headings and generate a table of contents.
type Document struct {
	Body  template.HTML
	Title string
}

var hairSpacer = regexp.MustCompile(`(\s*)(–|—)(\s*)`)

func useHairSpaces(content []byte) []byte {
	return hairSpacer.ReplaceAll(content, []byte("\u200A${2}\u200A"))
}

// GetHTMLPathname returns a pathname whose basename ends in .html and which
// resides in outputDirname, given an input pathname. If outputDirname is empty,
// returns a pathname that is a sibling of pathname. If the result is equal to
// pathname, returns an error.
func GetHTMLPathname(pathname, outputDirname string) (string, error) {
	i := strings.LastIndex(pathname, ".")
	if i == -1 {
		i = len(pathname)
	}
	var htmlPathname string
	if outputDirname != "" {
		htmlPathname = path.Join(outputDirname, pathname[:i]+".html")
	} else {
		htmlPathname = pathname[:i] + ".html"
	}
	if htmlPathname == pathname {
		return "", fmt.Errorf("cannot overwrite %q", pathname)
	}
	return htmlPathname, nil
}

// GenerateDocument executes template t with the document content, writing to w.
// It completes [Document.Title] by scanning content for the 1st <h1> tag.
func GenerateDocument(t *template.Template, content []byte, w io.Writer) error {
	//content = useHairSpaces(content)
	document := Document{Body: template.HTML(content)}
	parsed, e := htmlquery.Parse(strings.NewReader(string(content)))
	if e != nil {
		return e
	}
	title := htmlquery.FindOne(parsed, "//h1")
	if title == nil {
		return fmt.Errorf("missing h1")
	}
	document.Title = htmlquery.InnerText(title)
	if document.Title == "" {
		return fmt.Errorf("title has no text")
	}
	t.Execute(w, &document)
	return nil
}

// Minify minifies HTML, CSS, and JS in r and writes the result to w.
func Minify(w io.Writer, r io.Reader) {
	minifier.Minify("text/html", w, r)
}

// GenerateHTMLFile executes template t with the document contentPathname,
// writing the minified result to a new file in outputDirname. (See
// [GetHTMLPathname], [GenerateDocument], and [Minify].)
func GenerateHTMLFile(t *template.Template, contentPathname, outputDirname string) error {
	content, e := os.ReadFile(contentPathname)
	if e != nil {
		return e
	}

	htmlBytes := bytes.NewBuffer(make([]byte, 0, len(content)))
	if e := GenerateDocument(t, content, htmlBytes); e != nil {
		return e
	}

	minified := bytes.NewBuffer(make([]byte, 0, htmlBytes.Len()))
	Minify(minified, htmlBytes)

	htmlPathname, e := GetHTMLPathname(contentPathname, outputDirname)
	if e != nil {
		return e
	}
	if e := os.MkdirAll(path.Dir(htmlPathname), 0755); e != nil {
		return e
	}
	htmlFile, e := os.Create(htmlPathname)
	if e != nil {
		return e
	}
	if _, e := io.Copy(htmlFile, minified); e != nil {
		return e
	}
	return htmlFile.Close()
}

// GenerateHTML executes template t with the document contentPathname, writing
// the minified result to w. (See [GenerateDocument] and [Minify].)
func GenerateHTML(t *template.Template, contentPathname string, w io.Writer) error {
	content, e := os.ReadFile(contentPathname)
	if e != nil {
		return e
	}

	htmlBytes := bytes.NewBuffer(make([]byte, 0, len(content)))
	if e := GenerateDocument(t, content, htmlBytes); e != nil {
		return e
	}

	minified := bytes.NewBuffer(make([]byte, 0, htmlBytes.Len()))
	Minify(minified, htmlBytes)

	_, e = io.Copy(w, minified)
	return e
}
