define([
	'marionette', 'app/models/piece'
	],
	function(Marionette, Piece){
		var Rook = Piece.extend({
			initialize:function(){
				this.set('imgname', 'rook');
			}
		});
		return Rook;
	}
);