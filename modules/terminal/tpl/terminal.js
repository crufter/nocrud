$(function(){
	var uname = $("#logged-in-as").attr("data-user")
	var commands = []
	var ccounter = 0
	function setCmd(i) {
		$("#terminal-inp").val(commands[i])
		var l = commands[i].length
		setCaretToPos($("#terminal-inp")[0], l+1)
	}
	$(document).bind("keydown", function(e){
		switch (e.keyCode) {
		// Enter
		case 13:
			if (!e.shiftKey) {
				$("#display").append("<div class=\"command\">"+uname+" > "+$("#terminal-inp").val().replace(new RegExp("\n", 'g'), "<br />")+"</div>")
				$("#terminal-inp").submit()
				$("#terminal-inp").val("")
			} else {
				var where = doGetCaretPosition($("#terminal-inp")[0])
				var text = $("#terminal-inp").val()
				$("#terminal-inp").val(text.splice(where, 0, "\n"))
				setCaretToPos($("#terminal-inp")[0], where+1)
			}
			return false
		break
		// Up arrow
		case 38:
			if (ccounter == 0) {
				if (commands.length == 0) {
					return
				}
				var last = commands.length-1
				setCmd(last)
				ccounter = last
			} else {
				setCmd(ccounter-1)
				ccounter--
			}
		break
		// Down arrow
		case 40:
			if (ccounter == commands.length) {
				setCmd(0)
				ccounter = 0
			} else {
				setCmd(ccounter+1)
				ccounter++
			}
		break
		// V
		case 118:
			if (e.ctrlKey) {
				
			}
		}
	})
	
	$("#terminal-inp").focus()
	$(document).click(function(){
		//$("#terminal-inp").focus()
	})
	
	function cls() {
		$("#display").html("&nbsp;")
	}
	
	var builtins = {
		"cls": cls
	}
	
	$("form").on("submit", function() {
		var command = $("#terminal-inp").val()
		commands.push(command)
		ccounter = commands.length
		var cmd = command.split(" ")[0]
		if (builtins[cmd] != undefined) {
			builtins[cmd]()
			return
		}
		$.ajax({
			"url": 	"/terminal/execute",
			"data":	{
				"json":true,
				"1script": $.trim($("#terminal-inp").val())
			},
			"dataType": "json",
			"type": "POST",
			"success": function(data) {
				console.log(data)
				if (data["error"] == undefined) {
					var output = data["main"]
					output = output.replace(new RegExp("\n", 'g'), "<br />")
					$("#display").append("<div>"+output+"</div>")
				} else {
					$("#display").append("<div class=\"error\">"+data["error"]+"</div>")
				}
			},
			"error": function(d){
				$("#display").append("<div class=\"error\">http error "+d.responseText+"</div>")
			}
		})
		return false
	})
	
})

// from http://stackoverflow.com/questions/9370197/caret-position-cross-browser
function doGetCaretPosition (oField) {
	// Initialize
	var iCaretPos = 0
	// IE Support
	if (document.selection) { 
		// Set focus on the element
		oField.focus()
		// To get cursor position, get empty selection range
		var oSel = document.selection.createRange()
		// Move selection start to 0 position
		oSel.moveStart ('character', -oField.value.length)
		// The caret position is selection length
		iCaretPos = oSel.text.length;
	} else if (oField.selectionStart || oField.selectionStart == '0') {	// Firefox support
		iCaretPos = oField.selectionStart;
	}
	// Return results		// <- You don't say.
	return (iCaretPos);
}

function setSelectionRange(input, selectionStart, selectionEnd) {
	if (input.setSelectionRange) {
		input.focus()
		input.setSelectionRange(selectionStart, selectionEnd)
	}
	else if (input.createTextRange) {
		var range = input.createTextRange()
		range.collapse(true)
		range.moveEnd('character', selectionEnd)
		range.moveStart('character', selectionStart)
		range.select()
	}
}

function setCaretToPos(input, pos) {
	setSelectionRange(input, pos, pos)
}

// Ouch, prototype...
String.prototype.splice = function( idx, rem, s ) {
	return (this.slice(0,idx) + s + this.slice(idx + Math.abs(rem)))
};