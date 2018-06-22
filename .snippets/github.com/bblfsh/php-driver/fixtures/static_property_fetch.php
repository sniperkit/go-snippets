<?php

// property name variations
A::$b;
A::$$b;
A::${'b'};

// array access
A::$b['c'];
A::$b{'c'};
