{{require header.t}}

<h1>{{$.main_noun}} / New:</h1>
{{$f := form "insert"}}
<form action="/{{$f.ActionPath}}" method="POST" enctype="multipart/form-data">
	{{$f.HiddenString}}
	{{range .main}}
		{{.key}}<br />
		{{$hname := concat .type "TypeHandler"}}
		{{$h := hook $hname}}
		{{if $h.Has}}
			{{$ret := $h.Fire "KeyPrefix" $f.KeyPrefix "key" .key "value" .value}}
			{{html $ret}}
		{{else}}
			<input name="{{$f.KeyPrefix}}{{.key}}"/><br />
		{{end}}
		<br />
	{{end}}
	<input type="submit" />
</form>

{{require footer.t}}