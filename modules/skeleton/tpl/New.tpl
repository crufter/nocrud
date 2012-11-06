{{require header.t}}

<h1>{{$.main_noun}} / New:</h1>
{{$f := form "insert"}}
<form action="/{{$f.ActionPath}}" method="POST" enctype="multipart/form-data">
	{{$f.HiddenString}}
	{{range .main}}
		{{.key}}<br />
		<input name="{{$f.KeyPrefix}}{{.key}}"/><br />
		<br />
	{{end}}
	<input type="submit" />
</form>

{{require footer.t}}