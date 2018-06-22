define([
	'marionette', 'app/collections/piecesCol',
	'app/views/pieceView'
	],
	function(Marionette, PiecesCol, PieceView){
		var PiecesColView = Marionette.CollectionView.extend({
			itemView: PieceView,
			className: 'whargarbl'
		});
		return PiecesColView;
	}
);