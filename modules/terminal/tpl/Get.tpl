<!DOCTYPE html>
<html>
<head>
	<script src="/tpl/terminal/jquery.min.1.7.js"></script>
	<script src="/tpl/terminal/textinputs_jquery.js"></script>
	<script src="/tpl/terminal/terminal.js"></script>
	<link rel="stylesheet" type="text/css" href="/tpl/terminal/terminal.css" />
</head>

<body>
	<div id="display">
		&nbsp;
	</div>
	
	<div id="left">
		<div id="logged-in-as" data-user="{{._user.name}}">{{._user.name}} > </div>
	</div>
	
	<div id="right">
		<form id="terminal" action="/terminal/execute">
			<textarea spellcheck=false id="terminal-inp" class="terminal-inp"></textarea>
		</form>
	</div>
	<div class="clearfix"></div>
	<div id="whitespace">&nbsp;</div>
</body>
</html>