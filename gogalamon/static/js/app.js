// Connect to the server's WebSocket
var serverSock = new WebSocket("ws://" + window.location.host + "/sock/");

function welcome() {

	// Keypress listener
	var listener = new window.keypress.Listener();
	listener.register_many([
		{
			"keys"       : "j",
			"on_keyup"   : function(e) {
				dismissElement(document.getElementById("welcome_message"));
				listener.destroy();
				document.getElementById("captain_name_input").focus();
			}
		}
	]);

	var captainName;
	// Captain name input onkeydown event
	document.getElementById("captain_name_input").onkeydown = function(e) {
		// If the enter key is pressed
		if((e.keyCode || e.charCode) === 13) {
			// Get the input text
			var chatInputBox = document.getElementById("captain_name_input");

			if(chatInputBox.value == "") {
				return;
			} // end if

			captainName = chatInputBox.value;

			// Send the captain name
			serverSock.send(JSON.stringify({
				Event   : "username",
				User    : captainName
			}));

			dismissElement(document.getElementById("choose_captain_name"));

			document.getElementById("captain_name_input").blur();
		} // end if
	};

	// If the player chooses the Gopher team
	document.getElementById("choose_gophers").onclick = function () { 
		serverSock.send(JSON.stringify({
			Event   : "team",
			Team    : "gophers"
		}));

		dismissElement(document.getElementById("choose_team"));

		init();
	};
	// If the player chooses the Python team
	document.getElementById("choose_pythons").onclick = function () { 
		serverSock.send(JSON.stringify({
			Event   : "team",
			Team    : "pythons"
		}));

		dismissElement(document.getElementById("choose_team"));

		init();
	};


} // end welcome()


function dismissElement(elementToDismiss) {
	elementToDismiss.classList.add("done");
} // end dismissElement()


function init() {
	serverSock.onmessage = function(message) {
		var jsonMessage = JSON.parse(message.data);


		if(jsonMessage.Event == "chatMessage") {
			// Add the chat message to the output box
			var chatOutput = document.getElementById("chat_output");
			chatOutput.innerHTML += jsonMessage.Data.User + ": " + (jsonMessage.Data.Message).replace(/[<>]/g, '') + "<br>";

			// Scroll to bottom of textbox
			chatOutput.scrollTop = chatOutput.scrollHeight;
		} else if(jsonMessage.Event == "screenUpdate") {
			viewCenter.x = jsonMessage.Data.ViewX;
			viewCenter.y = jsonMessage.Data.ViewY;

			updateData = jsonMessage.Data;

			update(updateData);
		} else if(jsonMessage.Event == "ping") {
			serverSock.send(JSON.stringify({
				Event   : "pong"
			}));
		} else if(jsonMessage.Event == "sound") {
			playSound(jsonMessage.Data);
		}
	}; // end onmessage()


	// Init the stage
	var stage = new createjs.Stage("mainCanvas");
	var mainCanvas = document.getElementById("mainCanvas");

	// Init the mini map
	var miniMap = new createjs.Stage("miniMap");

	// Init the health bar
	var health = new createjs.Stage("health");
	var healthCanvas = document.getElementById("health");

	// Init capturing planet bar
	var capture = new createjs.Stage("capture");
	var captureCanvas = document.getElementById("capture");


	// Init the location, in map space, of the center (and therefor our player) of our view
	var viewCenter = {
		x : null,
		y : null
	}


	// Keypress listener
	var listener = new window.keypress.Listener();

	listener.register_many([
		{
			"keys"       : "w",
			"on_keydown" : function() {
				serverSock.send(JSON.stringify({
					Event: "w down"
				}));
			},
			"on_keyup"   : function(e) {
				serverSock.send(JSON.stringify({
					Event: "w up"
				}));
			}
		},
		{
			"keys"       : "a",
			"on_keydown" : function() {
				serverSock.send(JSON.stringify({
					Event: "a down"
				}));
			},
			"on_keyup"   : function(e) {
				serverSock.send(JSON.stringify({
					Event: "a up"
				}));
			}
		},
		{
			"keys"       : "s",
			"on_keydown" : function() {
				serverSock.send(JSON.stringify({
					Event: "s down"
				}));
			},
			"on_keyup"   : function(e) {
				serverSock.send(JSON.stringify({
					Event: "s up"
				}));
			}
		},
		{
			"keys"       : "d",
			"on_keydown" : function() {
				serverSock.send(JSON.stringify({
					Event: "d down"
				}));
			},
			"on_keyup"   : function(e) {
				serverSock.send(JSON.stringify({
					Event: "d up"
				}));
			}
		},
		{
			"keys"       : "j",
			"on_keydown" : function() {
				serverSock.send(JSON.stringify({
					Event: "f down"
				}));
			},
			"on_keyup"   : function(e) {
				serverSock.send(JSON.stringify({
					Event: "f up"
				}));
			}
		}
	]);


	// Get the chat input box
	var chatInput = document.getElementById('chat_input');
	// Stop listening for keyboard events for the canvas when the chat box is focussed
	chatInput.addEventListener("focus", chatInputFocussed);
	function chatInputFocussed() {
		listener.stop_listening();
	} // end chatInputFocussed()
	// Start listening again when it loses focus
	chatInput.addEventListener("blur", chatInputFocusLost);
	function chatInputFocusLost() {
		listener.listen();
	} // end chatInputFocusLost()


	// Text chat input onkeydown event
	document.getElementById("chat_input").onkeydown = function(e) {
		// If the enter key is pressed
		if((e.keyCode || e.charCode) === 13) {
			// Get the input text
			var chatInputBox = document.getElementById("chat_input");

			if(chatInputBox.value == "") {
				return;
			} // end if

			// Send the chat message
			serverSock.send(JSON.stringify({
				Event: "chatMessage",
				Message : chatInputBox.value
			}));

			// Add the chat message to the output box
			var chatOutput = document.getElementById("chat_output");
			chatOutput.innerHTML += "You: " + (chatInputBox.value).replace(/[<>]/g, "") + "<br>";

			// Scroll to bottom of textbox
			chatOutput.scrollTop = chatOutput.scrollHeight;

			// Reset the chat input box
			chatInputBox.value = "";
		} // end if
	};

	var currentNames = new Set();
	var nameCache = {};
	var ships = [];

	var sortFunction = function(obj1, obj2, options) {
		if (obj1.name > obj2.name) { return 1; }
		if (obj1.name < obj2.name) { return -1; }
		return 0;
	}


	// Set volume to half
	createjs.Sound.setVolume(0.5);
	// Register sounds
	createjs.Sound.registerSound("sounds/explosion0.wav", "explosion0");
	createjs.Sound.registerSound("sounds/explosion1.wav", "explosion1");
	createjs.Sound.registerSound("sounds/explosion2.wav", "explosion2");
	createjs.Sound.registerSound("sounds/laser0.wav", "laser0");
	createjs.Sound.registerSound("sounds/laser1.wav", "laser1");
	createjs.Sound.registerSound("sounds/laser2.wav", "laser2");
	createjs.Sound.registerSound("sounds/hit0.wav", "hit0");
	createjs.Sound.registerSound("sounds/hit1.wav", "hit1");
	createjs.Sound.registerSound("sounds/hit2.wav", "hit2");
	createjs.Sound.registerSound("sounds/shipThrust.wav", "shipThrust");

	function playSound(soundToPlay) {
		var instance = createjs.Sound.play(soundToPlay);
		if (soundToPlay == "laser0" || soundToPlay == "laser1"
		 || soundToPlay == "laser2") {
			instance.setVolume(0.25);
		}
	} // end playSound()

	var thrustingSound = createjs.Sound.play("shipThrust");

	function update(updateData) {
		// To cache an object: DisplayObject.cache()

		if (thrustingSound.playState === "playFailed"){
			thrustingSound = createjs.Sound.play("shipThrust");
			thrustingSound.play()
			thrustingSound.setLoop(100000000);
		} else {
			if (thrustingSound.paused && updateData.EngineSound) {
				thrustingSound.resume();
			} else if (!thrustingSound.paused && !updateData.EngineSound) {
				thrustingSound.pause();
			}
		}

		var newNames = new Set();
		for (var i = 0; i < updateData.Objs.length; i++){
			newNames.add(updateData.Objs[i].I)
		}
		newObjects = updateData.Objs;

		removeOldChildren(newNames);

		// Place the far starfield
		for (var i = mod(viewCenter.x * -0.1) - 512; i < mainCanvas.width; i += 512) {
			for (var j = mod(viewCenter.y * -0.1) - 512; j < mainCanvas.height; j += 512) {
				var starFieldFar = new createjs.Bitmap("img/starfield_far.png");
				starFieldFar.x = i;
				starFieldFar.y = j;
				starFieldFar.name = -3;

				stage.addChild(starFieldFar);
			};
		};

		// Place the mid starfield
		for (var i = mod(viewCenter.x * -0.4) - 512; i < mainCanvas.width; i += 512) {
			for (var j = mod(viewCenter.y * -0.4) - 512; j < mainCanvas.height; j += 512) {
				var starFieldNear = new createjs.Bitmap("img/starfield_near.png");
				starFieldNear.x = i;
				starFieldNear.y = j;
				starFieldNear.name = -2;

				stage.addChild(starFieldNear);
			};
		};

		// Place the near starfield
		for (var i = mod(viewCenter.x * -0.9) - 512; i < mainCanvas.width; i += 512) {
			for (var j = mod(viewCenter.y * -0.9) - 512; j < mainCanvas.height; j += 512) {
				var starFieldMid = new createjs.Bitmap("img/starfield_middle.png");
				starFieldMid.x = i;
				starFieldMid.y = j;
				starFieldMid.name = -1;

				stage.addChild(starFieldMid);
			};
		};


		// Create and place each new object we're sent
		for(var i = 0; i < updateData.Objs.length; i++) {
			// Get the object we want to render
			var currentObject = updateData.Objs[i];

			var objectBitmap;
			var addChildBool;

			if(!currentNames.has(currentObject.I)) {
				// Create the bitmap object
				objectBitmap = new createjs.Bitmap("img/" + currentObject.N + ".png");

				// Set the bitmap name to its unique id
				objectBitmap.name = currentObject.I;

				addChildBool = true;
				nameCache[currentObject.I] = objectBitmap;
			} else {
				objectBitmap = nameCache[currentObject.I];
				addChildBool = false;
			} // end if/else

			// Set the middle of the image
			objectBitmap.regX = objectBitmap.image.width / 2;
			objectBitmap.regY = objectBitmap.image.height / 2;

			objectBitmap.x = Math.round(currentObject.X - viewCenter.x + mainCanvas.width/2);
			objectBitmap.y = Math.round(currentObject.Y - viewCenter.y + mainCanvas.height/2);

			var currentRotation = objectBitmap.rotation;
			var targetRotation = currentObject.R;

			if(currentRotation < 5 && currentRotation > -5) {
				objectBitmap.rotation = 0;
			} else if(currentRotation < 95 && currentRotation > 85) {
				objectBitmap.rotation = 90;
			} else if(currentRotation < 185 && currentRotation > 175) {
				objectBitmap.rotation = 180;
			} else if(currentRotation < 275 && currentRotation > 265) {
				objectBitmap.rotation = 270;
			} 

			if(targetRotation < 0) {
				targetRotation += 360;
			} else if(targetRotation >= 360) {
				targetRotation -= 360;
			}

			if(currentObject.N.indexOf("ship") != -1) {
				createjs.Tween.get(objectBitmap, {override:true}).to({rotation:targetRotation}, 100, createjs.Ease.getPowInOut(2));
			} else {
				objectBitmap.rotation = currentObject.R;
			}

			// If the object is already on the stage, don't add it
			if(addChildBool) {
				stage.addChild(objectBitmap);
			} // end if
		} // end for

		currentNames = newNames;
		
		stage.sortChildren(sortFunction);

		stage.update();


		// Start updating mini map stuff

		var miniMapSize = {
			width  : document.getElementById("miniMap").width,
			height : document.getElementById("miniMap").height
		};

		miniMap.removeAllChildren();

		var circleGraphic = new createjs.Graphics().beginFill("Black").drawCircle(100, 100, 100);
		var circle = new createjs.Shape(circleGraphic);
		miniMap.addChild(circle);

		for(var i = 0; i < updateData.Planets.length; i++) {
			var currentPlanet = updateData.Planets[i];

			var planetBitmap = new createjs.Bitmap("img/planet_icon_" + currentPlanet.Team + ".png");

			planetBitmap.regX = planetBitmap.image.width / 2;
			planetBitmap.regY = planetBitmap.image.height / 2;

			planetBitmap.x = Math.round(currentPlanet.X/100 + 100);
			planetBitmap.y = Math.round(currentPlanet.Y/100 + 100);

			miniMap.addChild(planetBitmap);
		}

		if(updateData.Ships != null){
			ships = updateData.Ships;
		}

		for(var i = 0; i < ships.length; i++){
			var currentship = ships[i];

			var shipBitmap = new createjs.Bitmap("img/blip_" + currentship.Team + ".png");
			shipBitmap.regX = shipBitmap.image.width / 2;
			shipBitmap.regY = shipBitmap.image.height / 2;

			shipBitmap.x = Math.round(currentship.X/100 + 100);
			shipBitmap.y = Math.round(currentship.Y/100 + 100);

			miniMap.addChild(shipBitmap);
		}

		var miniShipBitmap = new createjs.Bitmap("img/marker.png");
		miniShipBitmap.regX = miniShipBitmap.image.width / 2;
		miniShipBitmap.regY = miniShipBitmap.image.height / 2;

		miniShipBitmap.x = Math.round(viewCenter.x/100 + 100);
		miniShipBitmap.y = Math.round(viewCenter.y/100 + 100);

		miniMap.addChild(miniShipBitmap);

		miniMap.update();


		// Health bar stuff

		health.removeAllChildren();

		for(var i = 1; i <= 10; i++) {
			if(i/10 <= updateData.Health) {
				var healthBar = new createjs.Bitmap("img/bar_health.png");
				healthBar.x = (i - 1) * healthBar.image.width;

				health.addChild(healthBar);
			} // end if
		} // end for

		health.update();


		// Capturing bar stuff

		capture.removeAllChildren();

		//console.log(updateData);

		if(updateData.PlanetAllegance != "") {
			for(var i = 1; i <= 10; i++) {
				if(i/10 <= updateData.AllegancePercent) {
					var captureBar = new createjs.Bitmap("img/bar_" + updateData.PlanetAllegance + ".png");
					captureBar.x = (i - 1) * captureBar.image.width;

					capture.addChild(captureBar);
				} // end if
			} // end for
		} // end if

		capture.update();

	} // end update()

	function removeOldChildren(newNames) {
		var toRemove = []
		for(var i = 0; i < stage.getNumChildren(); i++) {
			var child = stage.children[i];
			if (!newNames.has(child.name)){
				toRemove.push(child)
				delete nameCache[child.name]
			}
		}
		for(var i = 0; i < toRemove.length; i++){
			stage.removeChild(toRemove[i]);
		}
	} // end removeOldChildren()

	function mod(z) {
		z = z % 512;

		if(z < 0) {
			z += 512;
		} // end if

		return z;
	} // end mod()

} // end init()

