<?php

namespace app\admin\controller;

use think\Controller;
use think\Request;

class Index extends Controller
{
    public function index()
    {
        $article = model('Article');
        $admin = model('Admin')->count();
        $cate = model('Category')->count();
        return view()->assign([
            'art_num'=>$article->count(),
            'art' => $article
                ->order('create_time', 'desc')
                ->with('category')
                ->select(),
            'admin'=>$admin,
            'cate' =>$cate
        ]);
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
