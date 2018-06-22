function drawBase(color, x, y, size, name){
    var base = game.add.graphics(0, 0);
    base.beginFill(color, 1);
    base.drawCircle(x, y, size);

    //var text = "- phaser -\n with a sprinkle of \n pixi dust.";
    var style = { font: "32px Arial", fill: "#FFFFFF", align: "center" };
    var t = game.add.text(x, y+(size/2), name, style);

    if(name == playerName){
        timer = game.time.create(false);
        timer.loop(3000, requestDrawMinion, this);
        timer.start();
    }

    return base;
}

function requestDrawMinion(){
    ws.send(JSON.stringify({
        action:"request_minion",
        data:{
            name:playerName
        }
    }));
    console.log("Request minion build");
}
