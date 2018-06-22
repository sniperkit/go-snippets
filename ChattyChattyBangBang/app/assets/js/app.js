// Initialize app
var myApp = new Framework7();

// If we need to use custom DOM library, let's save it to $$ variable;
var $$ = Framework7.$;

// Initialize views
var mainView = myApp.addView('.view-main', {
	// Because we want to use dynamic navbar, we need to enable it for this view:
	dynamicNavbar: true
});

mainView.hideNavbar();


/* ==== Messasges Page ===== */
myApp.onPageInit('messages', function(page) {

	//myApp.showToolBar('.page[data-page="messages"');

	var conversationStarted = false;

	$$('.messagebar a.link').on('click', function() {
	
		var textarea = $$('.messagebar textarea');
		var messageText = textarea.val();
		
		if (messageText.length === 0) return;
		
		// empty textarea
		textarea.val('').trigger('change');
		textarea[0].focus();

		// Add Message
		myApp.addMessage({
			text: messageText,
			type: 'sent',
			day: !conversationStarted ? 'Today' : false,
			time: !conversationStarted ? (new Date()).getHours() + ':' + (new Date()).getMinutes() : false
		});

		conversationStarted = true;


	}); // end click

}); // end messages pageinit

myApp.onPageAfterAnimation('messages', function(page) {
	
	mainView.showToolbar();

});
