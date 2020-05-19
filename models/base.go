package models

// Resp is used for contruct uniform RESTful response
type Resp struct {
	Code *int        `json:"code"`
	Msg  *string     `json:"msg"`
	Data interface{} `json:"data",omitempty`
}

// Seg is a segment of an example
type Seg struct {
	Docs                            string `json:"docs,omitempty"`
	DocsRendered                    string `json:"docs_rendered,omitempty"`
	Code, CodeRendered, CodeForJs   string
	CodeEmpty, CodeLeading, CodeRun bool
}

// ExampleBase is example info extracted from gernerate.json
type ExampleBase struct {
	InnerID            *string  `json:"_id,omitempty" bson:"_id"`
	ID                 *string  `json:"id,omitempty" bson:"id"`
	NO                 *int64   `json:"no,omitempty" bson:"no"`
	Name               *string  `json:"name,omitempty" bson:"name"`
	DocsMarkup         *string  `json:"docs_markup,omitempty"`
	RunDocs            *string  `json:"run_docs,omitempty"`
	GoCode             *string  `json:"go_code,omitempty" bson:"gocode"`
	GoCodeHash         *string  `json:"gocodehash,omitempty" bson:"gocodehash"`
	URLHash            *string  `json:"url_hash,omitempty" bson:"urlhash"`
	Segs               [][]*Seg `json:"segs,omitempty" bson:"segs"`
	HighlightCode      *string  `json:"highlight_code"`
	HighlightCodeClean *string  `json:"highlight_code_clean"`
}
