package lesson

import (
	"log"
	"path/filepath"
	"testing"
)

func TestParseLesson(t *testing.T) {
	file := filepath.Join("test", "test.md")
	lesson, _, err := PasrseLesson(&file, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(string(lesson))
}
