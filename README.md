# content2html

Package content2html implements a very simple HTML templating system for
documents. It applies an `html/template` and fills it with the body text from an
input file, discovering the document’s `<title>` from its 1st `<h1>`. This way,
you can write plain HTML, and generate complete documents with templatized
boilerplate.
