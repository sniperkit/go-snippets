goog.provide('launcher.scriptpage');
goog.require('goog.dom');
goog.require('goog.events');
goog.require('goog.net.WebSocket');

launcher.scriptpage = function() {
    var ws = new goog.net.WebSocket();
    this.ws = ws;
    goog.events.listen(ws, goog.net.WebSocket.EventType.MESSAGE,
			     this.onMessage, undefined, this);
    try {
	ws.open('ws://'+goog.global['location']['hostname']+':8080/ws');
    } catch (e) {
	console.log(e);
    }
    this.addButtonListeners();
};
goog.addSingletonGetter(launcher.scriptpage);

launcher.scriptpage.prototype.addButtonListeners = function() {
    var button = goog.dom.getElementByClass('button');
    goog.events.listen(button, goog.events.EventType.CLICK,
			     this.buttonHandler, undefined, this);
    var clear = goog.dom.getElementByClass('clear'); 
    goog.events.listen(clear, goog.events.EventType.CLICK,
			     this.clearHandler, undefined, this);   
};

launcher.scriptpage.prototype.getWs = function() {
    return this.ws;
};

launcher.scriptpage.prototype.buttonHandler = function() {
    this.getWs().send(goog.global['scriptName']);
};

launcher.scriptpage.prototype.clearHandler = function() {
    var output = goog.dom.getElementByClass('output');
    goog.dom.setTextContent(output, '');
};

launcher.scriptpage.prototype.onMessage = function(e) {
    var output = goog.dom.getElementByClass('output');
    output.innerHTML += e.message;
    output.scrollTop = output.scrollHeight;
};
launcher.scriptpage.getInstance();
