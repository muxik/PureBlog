<?php

namespace app\admin\controller;

use app\common\validate\Admin;
use think\Controller;
use think\Request;

class Login extends Controller
{
    public function index()
    {
        return view();
    }

    public function login(Request $request)
    {
        $data = $request->post();
        $data['password'] = md5($data['password']);
        if ($data['password'] == 'd41d8cd98f00b204e9800998ecf8427e') $data['password'] = '';

        $validate = new Admin();
        if (!$validate->scene('login')->check($data)) {
            $this->error($validate->getError());
        }
        unset($data['code']);
        $admins = model("Admin")->where($data)->find();

        if (!$admins) $this->error('登陆失败');

        session('admin', $data['username']);
        $this->success('登陆成功', '/admin/index');
    }
}
