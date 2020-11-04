<?php

namespace app\index\controller;

use think\Controller;
use think\Request;

class Article extends Controller
{
    /**
     * @param Request $request
     * @param $id
     */
    public function index(Request $request, $id)
    {
        // 阅读量+1

        $articles = model('Article')->with(['admin'])->find($id);
        $articles->setInc('read');

        $articles->save();

        $categorys = model('Category')->select();
        $web = model('Web')->find();

        return view()->assign([
            'art' => $articles,
            'tag' => explode('|',$articles['tag']),
            'web' => $web,
            'categorys' => $categorys,
        ]);
    }
}
