define([
	'marionette', 'app/models/piece',
	'app/views/pieceView'
	],
	function(Marionette, Piece, PawnView){
		var Pawn = Piece.extend({
			initialize:function(){
				this.set('imgname', 'pawn');
			}
		});
		return Pawn;
	}
);