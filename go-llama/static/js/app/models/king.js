define([
	'marionette', 'app/models/piece'
	],
	function(Marionette, Piece){
		var King = Piece.extend({
			initialize:function(){
				this.set('imgname', 'king');
			}
		});
		return King;
	}
);