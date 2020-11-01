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

        return view()->assign([
            'articles' => $articles,
            'page' => [
                'last'=> $articles->lastPage(),
                'current'  => $articles->currentPage(),
            ],
            'categorys' => $categorys,
        ]);
    }

}
