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

	if (filter_var(rtrim($test, "\n"), FILTER_VALIDATE_EMAIL)) {
		echo "true\n";
	}
	else {
		echo "false\n";
	}
}
?>
