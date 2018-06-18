#!/usr/bin/env php
<?php

function getDecision($name, $game, $answer) {
  echo "$name : $game : $answer \n";
}

function main() {

  $my_game = "my game";
  $my_answer = "my answer";


  $getDecisionCaller = function($name) use ($my_game, $my_answer) {
    getDecision($name, $my_game, $my_answer);
  };

  // Dont need this line, can just use: $getDecisionCaller('company 0');
  getDecision('company 0', $my_game, $my_answer);

  $getDecisionCaller('company 1');
  $getDecisionCaller('company 2');

}

main();

?>

