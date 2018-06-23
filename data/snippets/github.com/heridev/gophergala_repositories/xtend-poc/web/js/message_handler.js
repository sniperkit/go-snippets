function handleWebSocketMessage(evt){
    console.log(evt);
    var eventData = JSON.parse(evt.data);
    if( eventData.action == "init" ){
        // initGame({
        //     player:playerName
        // });
        // removeRegistrationElements();
    }
    if( eventData.action === "render_base" ){
        initGame({
            player:playerName
        });
        removeRegistrationElements();
        setTimeout(function(){
            eventData.data.players.forEach(function(baseInfo){
                console.log("building base");
                drawBase(baseInfo.color, baseInfo.x, baseInfo.y, 100, baseInfo.name);
            });
        }, 500);
    }
    if( eventData.action === "render_minion" ){
        // eventData.data = {
        //     name: "player1",
        //     x: 10,
        //     y: 10,
        //     color: ""
        // }
        renderMinion(eventData.data)
    }
}

function removeRegistrationElements(){
    $('#player_registration').hide();
    $('#player_wait_message').hide();
}
