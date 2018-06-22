define([
	'marionette', 'app/models/piece'
	],
	function(Marionette, Piece){
		var Knight = Piece.extend({
			initialize:function(){
				this.set('imgname', 'knight');
			}
		});
		return Knight;
	}
);