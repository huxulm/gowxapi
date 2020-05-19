package api

import (
	"bytes"
	"log"
	"testing"
)

func TestHlightCode(t *testing.T) {
	var buf bytes.Buffer
	HlightCode(&buf, Code)
	log.Println(buf.String())
}
