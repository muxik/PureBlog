<?php

namespace app\http\middleware;

class CheckLogin
{
    public function handle($request, \Closure $next)
    {
        // cookie 登录
//        if (!cookie('admin_user') || session('?admin')){
//            return redirect('/admin/login');
//        }

        if (cookie('admin_user') || !empty(session('user.id', '', 'admin_'))) {
            return $next($request);
        }
        return redirect('/admin/login');
    }
}
