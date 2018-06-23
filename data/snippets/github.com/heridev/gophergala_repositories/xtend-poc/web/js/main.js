var playerName;

function onNameEntered(){

    playerName = document.getElementById("player_name").value;

    if (playerName === ""){
        return;
    }

    disablePlayerRegistration();
    popWaitingMessage(playerName);
    connect(playerName, handleWebSocketMessage);
    // handleWebSocketMessage({
    //     data:"{\"action\":\"init\"}"
    // });
    // var sampleBase = {
    //     action: "render_base",
    //     data: {
    //         players: [{
    //             name: "yourname",
    //             x: 11,
    //             y: 12,
    //             color: "0xfff"
    //         }]
    //     }
    // };
    // handleWebSocketMessage({
    //     data: JSON.stringify(sampleBase)
    // });
}

function disablePlayerRegistration(){
    $('#player_name').prop('disabled',true);
    $('#player_ok').prop('disabled',true);
}

function popWaitingMessage(playerName){
    $('#player_wait_message').text("Ok " + playerName + " Please wait...");
}
