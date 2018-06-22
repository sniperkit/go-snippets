define('wsHandler', ['jquery', 'underscore', 'backbone'], function($, _, Backbone) {

	// var wsUrl = 'ws://192.168.1.25:800/ws';
	var wsUrl = 'ws://' + window.location.hostname + ':' + window.location.port + '/ws';
	// alert(wsUrl);

	var wsHandler = {};

	_.extend(wsHandler, Backbone.Events);

	wsHandler.on('message', function(msg){
		handleIncomingMessage(msg);
	});

	wsHandler.on('error', function(err){
		console.log(err);
	});

	wsHandler.on('opened', function(msg){
		wsHandler.connected = true;
		
	});

	wsHandler.on('closed', function(msg){
		wsHandler.connected = false;
	});


	var handleIncomingMessage = function(msg){
		var data = JSON.parse(msg.data);

		switch(data.type){
			case 'authentication_response':
				wsHandler.trigger('authentication_response', data.response, data.user);
				if(data.user) wsHandler.user = data.user;
				break;

			case 'signup_response':
				wsHandler.trigger('signup_response', data.response, data.user);
				if(data.user) wsHandler.user = data.user;
				break;

			case 'game_request':
				wsHandler.trigger('game_request', data.opponent);
				break;

			case 'game_response_rejection':
				wsHandler.trigger('game_response_rejection', data.response);
				break;

			case 'game_move_update':
				wsHandler.trigger('game_move_update', data.game);
				break;

			case 'game_chat':
				wsHandler.trigger('game_chat', data.from, data.message_id);
				break;

			case 'game_over':
				wsHandler.trigger('game_over', data.game);
				break;

			case 'game_get_valid_moves_response':
				wsHandler.trigger('game_get_valid_moves_response', data.moves);
				break;


			default:
				//console.log('Unknown data type', data.type, data);
				break;
		}
	};


	wsHandler.authenticate = function(username, password){
		wsHandler.socket.sendJSON({type: 'authentication_request', username: username, user_token: password});
	};

	wsHandler.register = function(username, password, verseai){
		wsHandler.socket.sendJSON({type: 'signup_request', username: username, user_token: password, is_ai: false, verses_ai: verseai});
	};

	wsHandler.gameResponse = function(accept){
		wsHandler.socket.sendJSON({type: 'game_response', response: (accept)?'ok':'not-ok'});
	};

	wsHandler.moveRequest = function(move){
		wsHandler.socket.sendJSON({type: 'game_move_request', move: move});
	};

	wsHandler.sendChat = function(chatId){
		wsHandler.socket.sendJSON({type: 'game_chat_request', message_id: parseInt(chatId)});
	};

	wsHandler.getValidMoves = function(move){
		wsHandler.socket.sendJSON({type: 'game_get_valid_moves_request', location: move});
	};


	wsHandler.on('all', function(){
		console.log(arguments);
	});


	//Start the socket
	wsHandler.socket = new WebSocket(wsUrl);
	wsHandler.socket.sendJSON = function(message){ wsHandler.socket.send(JSON.stringify(message)) };
 
	wsHandler.socket.onopen = function(evt) { wsHandler.trigger('opened', evt); }; 
	wsHandler.socket.onclose = function(evt) { wsHandler.trigger('closed', evt); }; 
	wsHandler.socket.onmessage = function(evt) { wsHandler.trigger('message', evt); }; 
	wsHandler.socket.onerror = function(evt) { wsHandler.trigger('error', evt); };

	return wsHandler;
});

