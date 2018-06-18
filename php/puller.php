<?php

$repos = [
	// to pull, 
	// visit http://your-domain.com/puller.php?repo=my-app
	'my-app'         => '/opt/production/my-app', 
];

// make sure that your web server 
// user (eg. nginx, www-data, apache, etc) has necessary 
// permission to execute the command inside the directory.
// uncomment below code to check who's your web server user
// die(`echo $(whoami)`);

$repo    = isset($repos[$_GET['repo']]) ? $repos[$_GET['repo']] : null;

if(!$repo) 
	http_response_code(404) && die("Not found: '{$_GET['repo']}' 🤔");

// reset everything & pull the latest from repo's origin
$gitCmd  = 'git reset --hard && git pull origin -f';
$command = 'cd ' . $repo . ' && ' . $gitCmd;
$exec    = shell_exec($command);
$msg     = $exec ?: "Failed to execute command: `$gitCmd`";

echo $_GET['repo'] . ' → ' . $msg;

die();
