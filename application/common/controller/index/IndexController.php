<?php


namespace app\common\controller\index;


use think\Controller;

class IndexController extends Controller
{
    protected $web;
    protected $article;
    protected $category;
    protected $connect;
    protected $link;
    protected $comment;


    /**
     * 类初始化
     * @throws \think\exception\DbException
     */
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
            ->with(['scategory'])
            ->order('sort', 'asc');

        $this->connect = model('LinkModel')
            ->where([['state', '>', 0], ['type', '=', 1]])->select();

        $this->link = model('LinkModel')
            ->where([['state', '>', 0], ['type', '=', 0]])->select();


        $this->assign([
            'connect' => $this->connect,
            'link'  => $this->link,
            'web' => $this->web,
            'category' => $this->category->select(),
            'new_article' => $this->article->limit(5)->select(),
            'tag' => $this->getTags(),
            'article' => $this->article->paginate(5)
        ]);


        parent::initialize();
    }

    /**
     * 获取文章Tag
     * @return array
     */
    protected function getTags()
    {
        $tag = [];
        $this->tags = array_unique($tag);
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
        return $tag;
    }


    protected function getComments($type, $id){
        $this->comment = model('CommentModel');

        return $this->comment
            ->where($type, $id)
            ->order('create_time', 'desc')
            ->select();
    }

}