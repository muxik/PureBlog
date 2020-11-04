<?php

namespace app\index\controller;

class Index
{
    public function index()
    {


        $web =model('Web')->find();
        $categorys = model('Category')->select();


        // 分类查询
        if (input('cate')){
            $articles = model('Article')
                ->where('category_id', input('cate'))
                ->order('create_time', 'desc')
                ->paginate(10);
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

        $articles = model('Article')
            ->order('create_time', 'desc')
            ->paginate(10);

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
