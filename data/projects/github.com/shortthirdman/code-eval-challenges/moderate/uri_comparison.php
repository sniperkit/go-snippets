<?php
$fh = fopen($argv[1], "r");
while (true) {
  $test = fgets($fh);

  if ($test == false) {
    break;
  }
  elseif ($test == "\n") {
    continue;
  }

  $st = explode(';', rtrim($test, "\n"));
  $st[0] = preg_replace_callback("/%([a-fA-F\d][a-fA-F\d])/",
    create_function('$matches', 'return chr(base_convert($matches[0], 16, 10));'),
    $st[0]);
  $st[1] = preg_replace_callback("/%([a-fA-F\d][a-fA-F\d])/",
    create_function('$matches', 'return chr(base_convert($matches[0], 16, 10));'),
    $st[1]);
  $st1 = preg_split("/\/+/", $st[0], 3);
  $st2 = preg_split("/\/+/", $st[1], 3);
  $st1[0] = strtolower($st1[0]);
  $st2[0] = strtolower($st2[0]);
  $st1[1] = preg_replace("/:(80)?$/", "", strtolower($st1[1]));
  $st2[1] = preg_replace("/:(80)?$/", "", strtolower($st2[1]));

  if ($st1[0] == $st2[0] && $st1[1] == $st2[1] && $st1[2] == $st2[2]) {
    echo "True\n";
  }
  else {
    echo "False\n";
  }
}
?>
