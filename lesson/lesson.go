package lesson

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"time"

	"golang.org/x/tools/present"
)

// File defines the JSON form of a code file in a page.
type File struct {
	Name    string `json:"name" bson:"name"`
	Content string `json:"content" bson:"content"`
	Hash    string `json:"hash" bson:"hash"`
}

// Page defines the JSON form of a tour lesson page.
type Page struct {
	ID      *string     `json:"id" bson:"_id"`
	Title   string      `json:"title" bson:"title"`
	Content string      `json:"content" bson:"content"`
	Files   []File      `json:"files" bson:"files"`
	Lesson  interface{} `json:"lesson" bson:"lesson"`
}

// Lesson defines the JSON form of a tour lesson.
type Lesson struct {
	ID          *string `json:"id" bson:"_id"`
	Title       string  `json:"title" bson:"title"`
	Description string  `json:"description" bson:"description"`
	Pages       []Page  `json:"pages,omitempty" bson:"pages,omitempty"`
	// Category specify a lesson's category
	Category   *string    `json:"category,omitempty" bson:"category,omitempty"`
	Seq        *int64     `json:"seq,omitempty" bson:"seq,omitempty"`
	CreateTime *time.Time `json:"create_time,omitempty" bson:"create_time,omitempty"`
}

// PasrseLesson parse a specific lesson
func PasrseLesson(lessonPath *string, playground bool) ([]byte, *Lesson, error) {
	present.PlayEnabled = playground
	f, err := os.Open(*lessonPath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	doc, err := present.Parse(prepContent(f), *lessonPath, 0)
	if err != nil {
		return nil, nil, err
	}

	lesson := Lesson{
		nil,
		doc.Title,
		doc.Subtitle,
		make([]Page, len(doc.Sections)),
		nil,
		nil,
		nil,
	}
	// Load tmplete
	tmpl, err := present.Template().Parse(TmplString)
	if err != nil {
		return nil, nil, err
	}
	for i, sec := range doc.Sections {
		p := &lesson.Pages[i]
		w := new(bytes.Buffer)
		if err := sec.Render(w, tmpl); err == nil {
			p.Title = sec.Title
			p.Content = w.String()
			codes := findPlayCode(sec)
			p.Files = make([]File, len(codes))
			for i, c := range codes {
				f := &p.Files[i]
				f.Name = c.FileName
				f.Content = string(c.Raw)
				hash := sha1.Sum(c.Raw)
				f.Hash = base64.StdEncoding.EncodeToString(hash[:])
			}
		}
	}
	w := new(bytes.Buffer)
	if err := json.NewEncoder(w).Encode(&lesson); err != nil {
		return nil, nil, err
	}
	return w.Bytes(), &lesson, nil
}

// findPlayCode returns a slide with all the Code elements in the given
// Elem with Play set to true.
func findPlayCode(e present.Elem) []*present.Code {
	var r []*present.Code
	switch v := e.(type) {
	case present.Code:
		if v.Play {
			r = append(r, &v)
		}
	case present.Section:
		for _, s := range v.Elem {
			r = append(r, findPlayCode(s)...)
		}
	}
	return r
}

// prepContent for the local tour simply returns the content as-is.
var prepContent = func(r io.Reader) io.Reader { return r }
