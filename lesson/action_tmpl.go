package lesson

// TmplString is represent of action.tmpl file for `present`
const TmplString = `
{{/*{
	This is the action template.
	It determines how the formatting actions are rendered.
	*/}}
	
	{{define "section"}}
		<h2>{{.Title}}</h2>
		{{range .Elem}}{{elem $.Template .}}{{end}}
	{{end}}
	
	{{define "list"}}
		<ul>
		{{range .Bullet}}
			<li>{{style .}}</li>
		{{end}}
		</ul>
	{{end}}
	
	{{define "text"}}
		{{if .Pre}}
		<pre>{{range .Lines}}{{.}}{{end}}</pre>
		{{else}}
		<p>
			{{range $i, $l := .Lines}}{{if $i}}{{template "newline"}}
			{{end}}{{style $l}}{{end}}
		</p>
		{{end}}
	{{end}}
	
	{{define "code"}}
		{{if .Play}}
			{{/* playable code is not displayed in the slides */}}
		{{else}}
			<div>{{.Text}}</div>
		{{end}}
	{{end}}
	
	{{define "image"}}
		<img src="{{.URL}}"{{with .Height}} height="{{.}}"{{end}}{{with .Width}} width="{{.}}"{{end}}>
	{{end}}
	
	{{define "link"}}
		<p class="link"><a href="{{.URL}}" target="_blank">{{style .Label}}</a></p>
	{{end}}
	
	{{define "html"}}{{.HTML}}{{end}}
	
	{{define "newline"}}
	{{/* No automatic line break. Paragraphs are free-form. */}}
	{{end}}	
`
