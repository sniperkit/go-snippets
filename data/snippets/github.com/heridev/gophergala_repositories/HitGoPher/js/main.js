	// Declare global variable
    var playerName = "aorjoa";
    var dataReceive;
    // Connect WebSocket
	var ws = new WebSocket("ws://localhost:12345/start");

	ws.onopen = function (){
            ws.send(JSON.stringify({"Action":"newPlayer", "Player": playerName}));
	};
	ws.onmessage = function(msg) {
		dataReceive = JSON.parse(msg.data);
        $("img").attr("src","images/cupcake.png"); 
        $("#"+dataReceive.position).attr("src","images/gopher_in_cake.png"); 
	};

    $(document).ready(function(){
        $('img').on("click",function(){
            var clickedId = this.id;
            if(clickedId == dataReceive.position){
                $('#'+this.id).attr("src","images/gopher_break.png"); 
                    setTimeout(function() {
                        ws.send(JSON.stringify({"Action":"hit", "Position":clickedId, "Player": playerName}));
                    },300);
            }
        });
    });