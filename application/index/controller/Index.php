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

        $category = model('CategoryModel')
            ->with(['scategory'])
            ->select();

//        return json($category);

        $viewData = [
            'article' => $article,
            'category' => $category
        ];

        return view()->assign($viewData);
    }
}
