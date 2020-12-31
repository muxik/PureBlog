<?php

namespace app\http\middleware;

class Auth
{
    public function handle($request, \Closure $next)
    {
        session(['prefix' => 'admin_',]);
        $session_id = session('user.id');
        $id = $session_id != null ? $session_id : cookie('admin_user')['id'];

        $info = model('AdminModel')->where('id', $id)->find();
        if (2 !==$info['state'] ){
            return json(['code' => 0, 'msg' => '权限不足']);
        }

        return $next($request);
    }
}
