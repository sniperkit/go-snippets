define([
	'marionette',
	'app/models/piece',
	'app/models/pawn',
	'app/models/rook',
	'app/models/knight',
	'app/models/bishop',
	'app/models/queen',
	'app/models/king',
	],
	function(Marionette, Piece, Pawn, Rook, Knight, Bishop, Queen, King){
		var PiecesCol = Backbone.Collection.extend({
			model:Piece,
		});
		return PiecesCol;
	}
);