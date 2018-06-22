function veboGetMic(successFn, errFn) {
	navigator.getUserMedia = navigator.getUserMedia || navigator.webkitGetUserMedia || navigator.mozGetUserMedia;
	navigator.getUserMedia({video: false, audio: true}, successFn, errFn);
}

function socketConn(id, cons, stream, talkStatusFn, dcFn) {
	var socket = new WebSocket("ws://"+ window.location.hostname + ":" + window.location.port + "/handshake/"+id);
	socket.onerror = function(evt) {
		console.log("Failed to connect.");
	}
	socket.onopen = function(evt) {
		socket.send(cons);
		console.log("Ready!");
	}
	
	var peer = new Peer(id, {key: 'rik09ia3irdaemi'});
	peer.on('call', function(call) {
		console.log("Answer!");
	    call.answer(stream); // Answer the call with an A/V stream.
	    call.on('stream', talkStatusFn);
	});
	
	socket.onmessage = function(msg) {
		var call = peer.call(toString(msg.data), stream);
		talkStatusFn()
	    call.on('stream', talkStatusFn);
		call.on('disconnected', dcFn);
	}
}