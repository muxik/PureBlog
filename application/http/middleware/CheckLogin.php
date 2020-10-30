<?php

namespace app\http\middleware;

class CheckLogin
{
    public function handle($request, \Closure $next)
    {
        if (!session('?admin')){
            return redirect("/admin");
        }
        return $next($request);
    }
}
