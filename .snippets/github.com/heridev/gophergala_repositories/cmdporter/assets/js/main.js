$(document).ready(function () {
	$('#message').hide();

	$('.command').on('click', function () {
		// Loading
		button = this;
		$(button).children(".glyphicon").removeClass("glyphicon-play").addClass("glyphicon-refresh");

		// Send cmd request
		$.ajax({
			url: '/cmd',
			type: 'POST',
			dataType: 'json',
			data: JSON.stringify({command: $(button).data("command")}),
		})
		.done(function(data) {
			// Stop loading
			switch (data.error) {
				case null:
					$(button).children(".glyphicon").removeClass("glyphicon-refresh").addClass("glyphicon-ok");
					setTimeout(function() {
						$(button).children(".glyphicon").removeClass("glyphicon-ok").addClass("glyphicon-play");
					}, 5000);
					break;
				default:
					$(button).children(".glyphicon").removeClass("glyphicon-refresh").addClass("glyphicon-remove");
					$("#message").html('Command ' + $(button).data("command") + 'failed :' + data.error);
					$('#message').show();
					setTimeout(function() {
						$(button).children(".glyphicon").removeClass("glyphicon-remove").addClass("glyphicon-play");
					}, 5000);
					break;
			}
		})
		.fail(function(error) {
			// Stop loading
			$(button).children(".glyphicon").removeClass("glyphicon-refresh").addClass("glyphicon-remove");
			$("#message").html('Command ' + $(button).data("command") + 'failed :' + error);
			$('#message').show();
			console.log(error);
		})
		.always(function() {
			setTimeout(function() {
				$("#message").hide();
			}, 5000);
		});
	});
});