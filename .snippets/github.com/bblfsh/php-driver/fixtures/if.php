<?php

if      ($a) { $x = 1;}
elseif  ($b) {}
elseif  ($c) {}
else         {}

if ($a) {} // without else

if      ($a):
elseif  ($b):
elseif  ($c):
else        :
endif;

if ($a): endif; // without else
