<?php
$e = unpack('S',"\x01\x00");
if ($e[1] === 1) {
  echo "LittleEndian\n";
} else {
  echo "BigEndian\n";
}
?>
