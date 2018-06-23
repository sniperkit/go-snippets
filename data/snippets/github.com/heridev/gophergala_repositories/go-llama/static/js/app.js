define(
	['app/views/boardView', 'app/views/interactionWindowView',
	'wsHandler'], 
	function(Board, InteractionWindow, wsHandler){
		var app = new Marionette.Application();

		app.addRegions({
			board: '#board-div',
			interactionWindow: '#interaction-window'
		});

		app.addInitializer(function(){
			app.board.show(new Board());
			app.interactionWindow.show(InteractionWindow);
			$('td').droppable({
				// drop: function(){
					
				// },
				accept: '.piece',
				hoverClass:'highlight'
			});
		});

		app.start();

		return app;
	}
);

		