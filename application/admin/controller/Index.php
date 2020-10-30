<?php

namespace app\admin\controller;

use think\Controller;
use think\Request;

class Index extends Controller
{
    public function index()
    {

        return view();
    }

    /**
     * 退出登录
     */
    public function out()
    {
        session(null);
        return redirect('/admin');
    }
}
