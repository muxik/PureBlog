<?php

namespace app\index\controller;

use think\Controller;

class Index extends Controller
{
    public function index()
    {
        $web = model('WebModel')
            ->find();

        $article = model('ArticleModel')
            ->field('category_id,title,pic,read,description,tag,top,create_time,admin_id')
            ->with(['category', 'admin'])
            ->where([['state', '>', 0]])
            ->order('create_time', 'desc');


        $category = model('CategoryModel')
            ->where([['state', '>', 0]])
            ->with(['scategory'])
            ->select();

        $tag = [];
        $tag = array_unique($tag);
        foreach (
            model('ArticleModel')
                ->field('tag')
                ->where('state', '>', 0)
                ->select()
            as
            $key => $vo
        ) {
            foreach (explode(',', $vo['tag']) as $k => $v) {
                array_push($tag, $v);
            }
        }
//        return json($category);

        $viewData = [
            'article' => $article->paginate(5),
            'new_article' => $article->limit(5),
            'category' => $category,
            'web' => $web,
            'tag' => $tag
        ];

        return view()->assign($viewData);
    }


}
