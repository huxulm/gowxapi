package helper

import (
	"io"

	"github.com/alecthomas/chroma/quick"
)

// HlightCode2Html ouput hl code into full html
func HlightCode2Html(w io.Writer, source string) error {
	// lexer, style := "go", "solarized-dark"
	// lexer, style := "go", "dracula"
	// lexer, style := "go", "monokai"
	// lexer, style := "go", "emacs"
	lexer, style := "go", "github"
	// Highlight(w, source, lexer, style)
	if err := quick.Highlight(w, source, lexer, "html", style); err != nil {
		return err
	}
	return nil
}
