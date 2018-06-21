$.validator.setDefaults({
	errorPlacement: function(error, element) {
		$(element).closest("div").addClass("has-error");
	}
});
$().ready(function() {
	$("#productform").validate();
});
