<?php

namespace app\http\middleware;

class CheckLogin
{
    public function handle($request, \Closure $next)
    {

        if (cookie('admin_user') || !empty(session('user.id', '', 'admin_'))) {
            return $next($request);
        }
        return redirect('/admin/login');
    }
}
