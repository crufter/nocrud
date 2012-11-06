{{require header.t}}

<h1>{{$.main_noun}} / New:</h1>
{{$f := form "update"}}
<form action="/{{$f.ActionPath}}" method="POST" enctype="multipart/form-data">
	{{$f.HiddenString}}
	{{range .main}}
		{{.key}}<br />
		{{$key := .key}}
		{{if eq .type "file"}}
			<input type="file" name="{{.key}}" multiple="multiple"/><br />	<!-- !!! -->
			<br />
			{{range .value}}
				{{.}} <a href="/{{url "delete-file" "key" $key "file" .}}">Delete</a><br />
			{{end}}
		{{else}}
			<input name="{{$f.KeyPrefix}}{{.key}}" value="{{.value}}"/><br />
		{{end}}
		<br />
	{{end}}
	<input type="submit" />
</form>

{{require footer.t}}