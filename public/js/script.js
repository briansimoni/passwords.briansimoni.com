/**
 * Created by brian on 4/22/17.
 */
$(document).ready(function() {


	// disable the signup form submit button
	// It will be re-enabled when the user has met password requirements
	$("#create-account-button").attr('disabled', true);
	$('#signup-form').submit(function(form) {
		var firstPassword = document.getElementById('password').value;
		var confirmPassword = document.getElementById('confirm-password').value;
		if(firstPassword !== confirmPassword) {
			console.log(firstPassword + " " + confirmPassword);
			alert('passwords must match!');
			form.preventDefault();
			return false;
		}
	});

	var QueryString = function () {
		// This function is anonymous, is executed immediately and
		// the return value is assigned to QueryString!
		var query_string = {};
		var query = window.location.search.substring(1);
		var vars = query.split("&");
		for (var i=0;i<vars.length;i++) {
			var pair = vars[i].split("=");
			// If first entry with this name
			if (typeof query_string[pair[0]] === "undefined") {
				query_string[pair[0]] = decodeURIComponent(pair[1]);
				// If second entry with this name
			} else if (typeof query_string[pair[0]] === "string") {
				var arr = [ query_string[pair[0]],decodeURIComponent(pair[1]) ];
				query_string[pair[0]] = arr;
				// If third or later entry with this name
			} else {
				query_string[pair[0]].push(decodeURIComponent(pair[1]));
			}
		}
		return query_string;
	}();

	$('#login-button').click(function() {
		$('#pwd-container').toggle();
		$('#login-container').toggle();
		if(QueryString.status) {
			$('.alert').toggle();
		}
	});

	$('#signup-button').click(function() {
		$('#login-container').toggle();
		$('#pwd-container').toggle();
		if(QueryString.status) {
			$('.alert').toggle();
		}
	});


	if(QueryString.status) {
		var message = QueryString.status;
		message = message.replace(/\+/g, ' ');
		console.log(message);
		$('.alert').html(decodeURIComponent(message));
		$('.alert').toggle();
	}

});

