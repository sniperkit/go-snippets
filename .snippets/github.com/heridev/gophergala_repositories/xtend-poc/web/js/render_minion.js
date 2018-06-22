function drawMinion(info){
    console.log("Minion built");
    var minion = game.add.graphics(0, 0);
    minion.beginFill(info.color, 1);
    minion.drawRect(info.x+(info.size/2), info.y, 8, 8);
}
