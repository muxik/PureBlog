<?php

namespace app\index\controller;

use think\Controller;

class Index extends Controller
{
    public function index()
    {
        $article = model('ArticleModel')
            ->field('category_id,title,pic,read,description,tag,top,create_time,admin_id')
            ->with(['category', 'admin'])
            ->where('state', '>', 0)
            ->select();
//        return  json($article);

        return view()->assign([
            'article' => $article
        ]);
    }
}
