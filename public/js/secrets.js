/**
 * Created by brian on 5/2/17.
 */

$(document).ready(function() {
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

	if(QueryString.status === "success") {
		console.log(QueryString.status);
		$('#success-message').show();
	} else if(QueryString.status === "deleted") {
		$('#success-message').html("Deleted Successfully.");
		$('#success-message').show();
	} else if(QueryString.status === "updated") {
		$('#success-message').html("Updated Successfully.");
		$('#success-message').show();
	} else if(QueryString.status) {
		$('#other-message').html("That application already exists.");
		$('#other-message').show();
	}


	$('.glyphicon-pencil').each(function(thing, value) {
		console.log(value.getAttribute('data-application'));
		var application = value.getAttribute('data-application');
		$(this).click(function() {
			$('#editModal input[type="text"]')[0].placeholder = application;
			$('#editModal input[type="text"]')[0].value = application;
			$('#editModal input[type="hidden"]')[0].value = application;
		})
	});

	$('#generate-button-1').click(function() {
		var randomstring = Math.random().toString(36).slice(-8);
		$('#myModal input[type="text"]')[1].value = randomstring
	});

	$('#generate-button-2').click(function() {
		var randomstring = Math.random().toString(36).slice(-8);
		$('#editModal input[type="text"]')[1].value = randomstring
	});
});