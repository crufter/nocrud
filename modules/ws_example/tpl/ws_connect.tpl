{{require header.t}}

<script>
	var connection = new WebSocket('ws://127.0.0.1:6061/ws/whatever/ws-hello');
	connection.onopen = function(e){
		console.log("opened websocket connection")
	}
	connection.onmessage = function(e){
		var server_message = e.data;
		console.log("server said:", server_message);
	}
	connection.onclose = function(e){
		console.log("server closed connection.")
	}
</script>

{{require footer.t}}