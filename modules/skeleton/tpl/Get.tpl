{{require header.t}}

<h1>{{$.main_noun}} / Get:</h1>
{{if .main}}
	{{range .main}}
		<a href="/{{$.main_noun}}/{{._id}}">{{fallback .title .name ._id}}</a><br />
		<br />
	{{end}}
{{else}}
	No {{$.main_noun}} yet.
{{end}}

{{require footer.t}}