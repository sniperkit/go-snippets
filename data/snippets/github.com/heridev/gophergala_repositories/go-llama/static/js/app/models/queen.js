define([
	'marionette', 'app/models/piece'
	],
	function(Marionette, Piece){
		var Queen = Piece.extend({
			initialize:function(){
				this.set('imgname', 'queen');
			}
		});
		return Queen;
	}
);