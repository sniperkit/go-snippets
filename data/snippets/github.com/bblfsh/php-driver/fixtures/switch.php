<?php

switch ($a) {
    case 0:
        $x = 0;
        break;
    // Comment
    case 1;
    default:
        $y = 1;
}

// alternative syntax
switch ($a):
endswitch;

// leading semicolon
switch ($a) { ; }
switch ($a): ; endswitch;
