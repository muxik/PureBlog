<?php

namespace app\index\controller;

class Index
{
    public function index()
    {
        $articles = model('Article')
            ->order('create_time', 'asc')
            ->paginate(10);
        $categorys = model('Category')->select();
        $web =model('Web')->find();

        return view()->assign([
            'articles' => $articles,
            'page' => [
                'last'=> $articles->lastPage(),
                'current'  => $articles->currentPage(),
            ],
            'web'  => $web,
            'categorys' => $categorys,
        ]);
    }

}
