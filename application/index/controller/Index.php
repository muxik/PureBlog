<?php

namespace app\index\controller;

class Index
{
    public function index()
    {
        $articles = model('Article')->select();
        $categorys = model('Category')->select();
        return view()->assign([
            'articles' => $articles,
            'categorys' => $categorys,
        ]);
    }
}
