<?php 
/**
 * Overload app's dotenv based on domain.
 * The file that contains domain-based environment 
 * is stored in the same directory of this file and has
 * same file name as the domain itself.
 * 
 * If this file is stored in `/domain_environments/` and you 
 * want to load domain env for `me.com`, store environment 
 * vars of `me.com` in `/domain_environments/me.com` file.
 */

function getDomainEnvFromArgv($removeFromArgv = false)
{
	$domainEnvPrefix = '-domain_env=';
	$args = $_SERVER['argv'];
	foreach ($args as $key => $arg) {
		if (stripos($arg, $domainEnvPrefix) !== false) {
			list(, $domain_env) = explode($domainEnvPrefix, $arg);
			
			if ($removeFromArgv) {
				unset($_SERVER['argv'][$key]);
			}

			return $domain_env;
		}
	}
	return null;
}

$domain_env = !empty($_SERVER['HTTP_HOST']) ? $_SERVER['HTTP_HOST'] : getDomainEnvFromArgv();

// overload .env file if `domain_env` is succesfully retrieved
if ($domain_env && file_exists(__DIR__ . '/' . $domain_env)) {
    $dotenv = new Dotenv\Dotenv(__DIR__, $domain_env);
    $dotenv->overload();
}
