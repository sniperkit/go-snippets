define([
	'marionette',
	'app/models/piece',
	'text!templates/piece.html',
	'wsHandler'
	],
	function(Marionette, Piece, PieceTemplate, wsHandler){
		var PieceView = Marionette.ItemView.extend({
			className:'piece',
			template:_.template(PieceTemplate),
			initialize:function(){
				_.bindAll(this, 'revert', 'startDrag');
			},
			onRender:function(){
				var color = (this.model.get('color') === window.whichColor);
				var user = (wsHandler.user.username === window.WhosMove);
				if(user && color){
					this.$el.draggable({
						snap: 'td',
						snapModel:'inner',
						opacity: 0.8,
						distance: 10,
						revert:this.revert,
						start:this.startDrag,
						stop:function(){
							$('.validSquare').each(function(index, element){
								$(element).removeClass('validSquare');
							});
						}
					});
				}

				this.$el.css('position','absolute');
				loc = this.model.get('location');
				this.$el.offset({left: loc[0] * 50, top: ((9 - loc[1]) * 50) - 30});
			},
			startDrag:function(event, ui){
				// console.log(event);
				var startLoc = this.model.get('location');
				var locChr = String.fromCharCode(96 + parseInt(startLoc[0]));
				var locationString = locChr + startLoc[1];
				if (window.YourColor === 'black'){
					locationString = locChr + (9 - startLoc[1]);
				}
				wsHandler.getValidMoves(locationString);
			},
			revert:function(socketObj){
				if(socketObj){
					// Drag success - this would be where we trigger submitting the move
					var dataset = socketObj[0].dataset;
					var newrow = dataset.row;
					var newcol = dataset.col;
					var loc = [newcol, newrow];
					var oldloc = this.model.get('location');
					var oldcol = oldloc[0];

					var oldchr = String.fromCharCode(96 + parseInt(oldcol));
					var newchr = String.fromCharCode(96 + parseInt(newcol));

					this.model.set('location', loc);

					var moveString = oldchr + oldloc[1] + '-' + newchr + loc[1];
					if(window.YourColor === 'black'){
						moveString = oldchr + (9 - oldloc[1]) + '-' + newchr + (9 - loc[1]);	
					}
					// $('#movelog').append(moveString + '<br>');

					wsHandler.moveRequest(moveString);

					// console.log('success!');
					return false;
				}
				else {
					// drag fail
					// console.log('fail!');
					return true;
				}
			}
		});
		return PieceView;
	}
);