/**
 * Created by brian on 4/22/17.
 */
$(document).ready(function() {
	$('#login-button').click(function() {
		$('#pwd-container').toggle();
		$('#login-container').toggle();
	});

	$('#signup-button').click(function() {
		$('#login-container').toggle();
		$('#pwd-container').toggle();
	});
});