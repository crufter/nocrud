{{require header.t}}

<h1>{{$.main_noun}} / LoginForm:</h1>

{{if logged_in}}
	{{$f := form "logout"}}
	<a href="/{{$f.ActionPath}}">Logout</a>
{{else}}
	{{$f := form "login"}}
	<form action="/{{$f.ActionPath}}" method="POST">
		name:<br />
		<input name="{{$f.KeyPrefix}}name" value=""/><br />
		<br />
		password:<br />
		<input type="password" name="{{$f.KeyPrefix}}password" value=""/><br />
		<br />
		<input type="submit">
	</form>
{{end}}

{{require footer.t}}