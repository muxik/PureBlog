<?php

namespace app\admin\controller;

use app\common\validate\Admin;
use think\App;
use think\Controller;
use think\Request;

class Login extends Controller
{

    public function index()
    {
        // 如果以session存在就跳转到主页
        if (session("?admin")) return redirect('/admin/index');
        return view();
    }

    public function login(Request $request)
    {
        // 如果以session存在就跳转到主页
        if (session("?admin")) return redirect('/admin/index');

        // 获取前台提交的数据进行调整
        $data = $request->post();
        $data['password'] = md5($data['password']);

        // 空字符串也可以被md5加密
        if ($data['password'] == 'd41d8cd98f00b204e9800998ecf8427e') $data['password'] = '';

        // 验证数据
        $validate = new Admin();
        if (!$validate->scene('login')->check($data)) {
            $this->error($validate->getError());
        }
        unset($data['code']);
        $admin = model("Admin")->where($data)->find();

        if (!$admin) $this->error('登录失败，用户名或密码错误');
        if ($admin['status'] != 1) $this->error('此账户已被管理员禁用');


        // 保存session
        $seData = [
            'username' => $admin['username'],
            'super' => $admin['super'] == 1,
            'id'    => $admin['id'],
        ];
        session('admin', $seData);
        $this->success('登陆成功', '/admin/index');
    }
}
