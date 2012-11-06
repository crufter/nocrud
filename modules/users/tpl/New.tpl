{{require header.t}}

<h1>{{$.main_noun}} / LoginForm:</h1>

{{$f := form "insert"}}
<form action="/{{$f.ActionPath}}" method="POST">
	name:<br />
	<input name="{{$f.KeyPrefix}}name" value=""/><br />
	<br />
	password:<br />
	<input type="password" name="{{$f.KeyPrefix}}password" value=""/><br />
	<br />
	<input type="submit">
</form>

{{require footer.t}}