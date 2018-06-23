define([
	'marionette',
	'wsHandler',
	'jquery',
	'text!templates/interactionWindow/loading.html',
	'text!templates/interactionWindow/loginForm.html',
	'text!templates/interactionWindow/registerForm.html',
	'text!templates/interactionWindow/gameRequest.html',
	'text!templates/interactionWindow/generic.html',
	'text!templates/interactionWindow/chat.html'

	],
	function(Marionette, wsHandler, $, loadingTpl, loginTpl, registerTpl, gameRequest, genericTpl, chatTpl){

		var chatMessages = {
			1: "Hello",
			2: "I am the King.",
			3: "You took my piece!",
			4: "Bishop! I need a bishop!",
			5: "I didn't like that Knight anyway...",
			6: ":)",
			7: ":(",
			8: "D:",
			9: ":D",
			10: "I'm gonna get you back for that.",
			11: "You sunk my battleship...",
			12: "Prepare for defeat!",
			13: "Nicely done.",
			14: "This isn't a game, this is a slaughter!",
			15: "Don't worry be happy :-)"
		};


		var dateNow = function(){
			var today = new Date();
		    var dd = today.getDate();
		    var mm = today.getMonth()+1; //January is 0!

		    var hh = today.getHours()+1;
		    var mm = today.getMinutes();
		    var ss = today.getSeconds();

		    var yyyy = today.getFullYear();
		    if(dd<10){
		        dd='0'+dd
		    } 
		    if(mm<10){
		        mm='0'+mm
		    } 

		    if(hh<10){
		    	hh = '0'+hh;
		    }

		    if(mm<10){
		    	mm = '0'+mm;
		    }

		    if(ss<10){
		    	ss = '0'+ss;
		    }


		    return dd+'/'+mm+'/'+yyyy+' '+hh+':'+mm+':'+ss;
		}


		var tplData = {msg: 'Connecting to server...'};

		var changeView = function(tpl, newData){
			tplData = newData || {};

			iw.template = _.template(tpl);
			iw.render();
		}


		var InteractionWindow = Marionette.Layout.extend({
			template: _.template(loadingTpl),

			render: function() {
				this.$el.html(this.template(tplData));
				return this;
			},

			events: {
				'click #login-button': 'doLogin',
				'click #i-want-to-register': 'showRegisterView',
				'click #i-want-to-login': 'showLoginView',
				'click #register-button': 'doRegister',
				'click #game-accept': 'acceptGame',
				'click #game-deny': 'denyGame',
				'click #chat-button': 'doChat'
			},
			

			doLogin: function(e){
				e.preventDefault();
				
				var username = $('#login-form input[name="username"]').val();
				var password = $('#login-form input[name="password"]').val();
				
				changeView(loadingTpl);

				return wsHandler.authenticate(username, password);
			},

			showRegisterView: function(){
				return changeView(registerTpl);
			},

			showLoginView: function(){
				return changeView(loginTpl);
			},

			doRegister: function(e){
				e.preventDefault();
				
				var username = $('#register-form input[name="username"]').val();
				var password = $('#register-form input[name="password"]').val();
				var ai = $('#register-form input[name="verseai"]').is(':checked');
				
				changeView(loadingTpl);

				return wsHandler.register(username, password, ai);
			},

			acceptGame: function(e){
				wsHandler.gameResponse(true);
				return changeView(loadingTpl, {msg: 'Waiting for other player to accept or deny'});
			},

			denyGame: function(e){
				wsHandler.gameResponse(false);
				return changeView(loadingTpl, {msg: 'Waiting for game'});
			},

			doChat: function(e){
				e.preventDefault();
				wsHandler.sendChat($('#chat-select').val());
			}

		});

		var iw = new InteractionWindow();

		var showLogin = function(){ changeView(loginTpl); };


		//If we haven't connected yet, we need to wait for the opened event
		//If we have, just show the login.
		if(wsHandler.connected){
			showLogin();
		}else{
			wsHandler.on('opened', showLogin);
		}


		wsHandler.on('authentication_response', function(code, user){
			if(code != 'ok'){
				return changeView(loginTpl, {error: 'Login incorrect'});
			}

			return changeView(loadingTpl, {msg: 'Waiting for game'});

		});

		wsHandler.on('signup_response', function(code, user){
			if(code != 'ok'){
				return changeView(registerTpl, {error: 'Registration incorrect'});
			}

			return changeView(loadingTpl, {msg: 'Waiting for game'});

		});


		wsHandler.on('game_request', function(opponent){
			//Only do this once if the game starts
			wsHandler.once('game_move_update', function(opponent){
				return changeView(chatTpl);
			});

			return changeView(gameRequest, {opponent: opponent});
		});

		wsHandler.on('game_response_rejection', function(opponent){
			return changeView(loadingTpl, {msg: 'Last game did not start.<br>Waiting for game'});
		});

		wsHandler.on('game_chat', function(from, messageId){
			$("#chat-history").append('<p>'+dateNow()+' - <b>'+from.username+'</b>: '+chatMessages[messageId]+'</p>');
			$("#chat-history").scrollTop($('#chat-history').prop("scrollHeight"));
		});

		wsHandler.on('game_over', function(game){
			return changeView(loadingTpl, {msg: 'Waiting for game'});
		});


		return iw;
	}
);