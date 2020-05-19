package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/jackdon/gowxapi/models"
	"github.com/julienschmidt/httprouter"
)

// Highlight rewrite quick.Highlight
func Highlight(w io.Writer, source string, lexer, style string) error {
	// Determine lexer.
	l := lexers.Get(lexer)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Determine formatter.
	f := html.New(html.Standalone(false), html.WithClasses(true))

	// Determine style.
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return f.Format(w, s, it)
}

// HlightCode is ...
func HlightCode(w io.Writer, source string) {
	lexer, style := "go", "github"
	Highlight(w, source, lexer, style)
	// quick.Highlight(w, source, lexer, "html", style)
}

const Code = `
var buf bytes.Buffer
fmt.Fprintf(&buf, "Size: %d MB.", 85)
s := buf.String()) // s == "Size: 85 MB."
`

// HighLight is ...
func HighLight(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	var buf bytes.Buffer
	HlightCode(&buf, string(body))
	msg := "ok"
	code := 0
	result := &models.Resp{Code: &code, Msg: &msg}
	result.Data = buf.String()
	resultB, _ := json.Marshal(result)
	fmt.Fprintf(w, fmt.Sprint(string(resultB)))
}
