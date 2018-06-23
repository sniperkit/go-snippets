#!/usr/bin/php
<?php

function run ($expr) {

  $cmd = "";

  if ( strpos($expr, ".scala") > -1) {
    $expr = explode('.', $expr);
    $expr = $expr[0];
  }

  if ( strpos($expr,  '/') > -1) {
    $expr = array_pop(explode('/',  $expr));
  }

  $cmd = "scala -classpath classes $expr";

  echo("Running: $cmd\n");
  system( $cmd );
}

if( count( $argv) == 1 )
  echo("USAGE: $argv[0] CLASSNAME\n");
else
  run( $argv[1]);
?>
