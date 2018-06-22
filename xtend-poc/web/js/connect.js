var ws;
function connect(playerName, onMessageHandler){
    if (!"WebSocket" in window)
    {
        console.log("WebSocket is not supported by your Browser!");
        return;
    }

    console.log("WebSocket is supported by your Browser!");
    ws = new WebSocket("ws://localhost:12345/api/start");

    ws.onopen = function()
    {
        ws.send(JSON.stringify({
            action:"new",
            data:{
                name:playerName
            }
        }));
        console.log("Message is sent...");
    };

    ws.onmessage = function(evt){
        onMessageHandler(evt);
    };

    ws.onclose = function()
    {
        alert("Connection is closed...");
    };
}
