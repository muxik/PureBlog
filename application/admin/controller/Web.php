<?php

namespace app\admin\controller;

use think\Controller;
use think\Request;

class Web extends Controller
{
    /**
     * 主页
     * @return \think\response\View
     */
    public function index()
    {
        $web = model('Web')->find();

        return view()->assign(['web' => $web]);
    }

    /**
     * 更新网站数据
     * @param Request $request
     */
    public function update(Request $request)
    {

    }
}
