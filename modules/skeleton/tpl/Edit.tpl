{{require header.t}}

<h1>{{$.main_noun}} / New:</h1>
{{$f := form "update"}}
<form action="/{{$f.ActionPath}}" method="POST">
	{{$f.HiddenString}}
	{{range .main}}
		{{.key}}<br />
		{{$hname := concat .type "TypeHandler"}}
		{{$h := hook $hname}}
		{{if $h.Has}}
			{{$ret := $h.Fire "KeyPrefix" $f.KeyPrefix "key" .key "value" .value}}
			{{html $ret}}
		{{else}}
			<input name="{{$f.KeyPrefix}}{{.key}}" value="{{.value}}"/><br />
		{{end}}
		<br />
	{{end}}
	<input type="submit" />
</form>

{{require footer.t}}