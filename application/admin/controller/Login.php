<?php

namespace app\admin\controller;

use think\Controller;
use think\Request;

class Login extends Controller
{
    public function index()
    {
        if (cookie('admin_user') || !empty(session('user.id', '', 'admin_'))) {
            $this->redirect('/admin/index');
        }

        return view();
    }

    public function login(Request $request)
    {
        $data = $request->only(['captcha', 'username', 'password', 'rememberMe' => null]);
        $result = model('AdminModel')->login($data);

        if (1 !== $result) {
            $this->error($result);
        }

        $this->success('登录成功！','/admin/index');
    }
}
