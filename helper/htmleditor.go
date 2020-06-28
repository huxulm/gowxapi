package helper

import (
	"bytes"

	"golang.org/x/net/html"
)

// DefaultCSS gives a default css
var DefaultCSS = []byte(`<meta charset="UTF-8"><style type="text/css">
html { background: grey; padding:0; margin:0; font-size: 16px; }
body{ padding:25px; margin:0; width: 380px; min-height:580px; background: whitesmoke; display:block; overflow:hidden;
}
</style>`)

// InjectStyle inject css into <style> tag
func InjectStyle(srcHtml, css []byte) []byte {
	r := bytes.NewReader(srcHtml)
	doc, err := html.Parse(r)
	if err != nil {
		return nil
	}
	ns, err := html.Parse(bytes.NewReader(css))
	if err != nil {
		return nil
	}
	n := ns.FirstChild.FirstChild.FirstChild
	n.Parent = nil
	n.PrevSibling = nil
	n.NextSibling = nil
	doc.FirstChild.FirstChild.AppendChild(n)

	w := new(bytes.Buffer)
	if err := html.Render(w, doc); err != nil {
		return nil
	}
	return w.Bytes()
}
