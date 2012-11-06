{{require header.t}}

<h1>{{$.main_noun}} / New:</h1>
{{$f := form "update"}}
<form action="/{{$f.ActionPath}}" method="POST">
	{{$f.HiddenString}}
	<input type="submit">
	<textarea id="code" name="{{$f.KeyPrefix}}json" style="display: block; width: 100%; height: 92%">{{.main}}</textarea>
</form>

{{require footer.t}}