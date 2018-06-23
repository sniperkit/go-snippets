#!/usr/bin/php
<?php

function compile($expr) {
  echo "Preparing ...\n";

  system("rm -rf classes");
  system("mkdir classes");

  $cmd = "";

  if($expr == "*") {
    $cmd = "scalac -d classes src/*.scala";
  } else if( strpos( $expr, ".scala") > -1) {
    $cmd = "scalac -d classes $expr";
  } else {
    $cmd = "scalac -d classes $expr.scala";
  }

  $cmd .= " -feature";

  echo "Compiling $cmd\n";
  system( $cmd );

  echo "Done.\n";
}

if( count($argv) == 1) {
  echo "No file specified. Using *\n";
  compile("*");
} else
  compile ($argv[1]);
?>
