<?php


namespace app\index\controller;


use think\Controller;

class IndexController extends Controller
{
    protected $web;
    protected $article;
    protected $category;
    protected $connect;


    protected function initialize()
    {
        $this->web = model('WebModel')
            ->find();

        $this->article = model('ArticleModel')
            ->field('id,category_id,title,pic,read,description,tag,top,create_time,admin_id')
            ->with(['category', 'admin'])
            ->where([['state', '>', 0]])
            ->order('create_time', 'desc');


        $this->category = model('CategoryModel')
            ->where([['state', '>', 0]])
            ->with(['scategory']);

        $this->connect = model('LinkModel')
            ->where([['state', '>', 0], ['type', '=', 1]])->select();

        $this->assign(['connect' => $this->connect]);



        parent::initialize();
    }
}