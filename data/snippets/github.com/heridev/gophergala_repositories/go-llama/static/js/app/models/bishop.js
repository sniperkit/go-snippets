define([
	'marionette', 'app/models/piece'
	],
	function(Marionette, Piece){
		var Bishop = Piece.extend({
			initialize:function(){
				this.set('imgname', 'bishop');
			}
		});
		return Bishop;
	}
);