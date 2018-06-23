var game;
var gameInfo;
var bases;

function initGame(gameInfo){
    game = new Phaser.Game(800, 600, Phaser.CANVAS, 'phaser-example', {
        create: create,
        update: update
    });

    function create(){

    }

    function update(){
        if (game.input.activePointer.isDown)
        {
            console.log(game.input.activePointer);
        }
    }
}
