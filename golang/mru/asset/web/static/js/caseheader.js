$('.list-group-item').on('mouseover', function(event) {
	event.preventDefault();
	$(this).closest('li').addClass('open');
});
$('.list-group-item').on('mouseout', function(event) {
	event.preventDefault();
	$(this).closest('li').removeClass('open');
});

$(".ignore").click(function() {
	return false
});
