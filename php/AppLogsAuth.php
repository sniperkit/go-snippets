<?php

namespace App\Http\Middleware;

use Auth;
use Closure;

class AppLogsAuth extends Authenticate
{
  /**
   * Handle an incoming request.
   *
   * @param  \Illuminate\Http\Request  $request
   * @param  \Closure  $next
   * @return mixed
   */
  public function handle($request, Closure $next)
  {
  	$envs = [
  	    'production'
  	];

  	if(in_array(app()->environment(), $envs)) {
  		$user = config('app-logs.user');
  		$password = config('app-logs.password');

  	    if($request->getUser() != $user && $request->getPassword() != $password) {
  	        $headers = ['WWW-Authenticate' => 'Basic'];
  	        return response('Unauthorized', 401, $headers);
  	    }
  	}

  	return $next($request);
  }
}
