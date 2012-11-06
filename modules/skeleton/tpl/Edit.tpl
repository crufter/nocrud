{{require header.t}}

<h1>{{$.main_noun}} / New:</h1>
{{$f := form "update"}}
<form action="/{{$f.ActionPath}}" method="POST">
	{{$f.HiddenString}}
	{{range .main}}
		{{.key}}<br />
		<input name="{{$f.KeyPrefix}}{{.key}}" value="{{.value}}"/><br />
		<br />
	{{end}}
	<input type="submit" />
</form>

{{require footer.t}}