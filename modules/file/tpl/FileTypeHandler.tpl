<input type="file" name="{{.key}}" multiple="multiple"/><br />	<!-- !!! -->
<br />
{{$key := .key}}
{{range .value}}
	{{.}} <a href="/{{url "delete-file" "key" $key "file" .}}">Delete</a><br />
{{end}}