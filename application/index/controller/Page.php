<?php


namespace app\index\controller;


use app\common\controller\index\IndexController;
use think\Request;

class Page extends IndexController
{

    protected $page;

    protected function initialize()
    {
        parent::initialize();
    }

    /**
     * 文章详情页
     * @param $id
     * @return \think\response\View
     * @throws \think\Exception
     */
    public function info($id)
    {

        $info = model('ArticleModel')
            ->field('source,id,category_id,content,title,pic,read,description,tag,top,create_time,admin_id')
            ->find($id);


        $comment = $this->getComments('article_id', $id);
        $info->setInc('read');

        $this->assign([
            'info' => $info,
            'comment' => $comment
        ]);


        return view('info');
    }

    /**
     * 关于
     * @return \think\response\View
     */
    public function about()
    {
        $comment = $this->getComments('page_id', 1);

        $this->assign(['comment' => $comment]);
        return view('page/about');
    }

    /**
     * 友人
     * @return \think\response\View
     */
    public function link()
    {

        $comment = $this->getComments('page_id', 2);
        $this->assign(['comment' => $comment]);
        return view('page/link');
    }

    /**
     * 碎碎念
     * @return \think\response\View
     */
    public function timeLine()
    {
        $comment = $this->getComments('page_id', 3);

        $shuo = model('ShuoModel')
            ->order('create_time', 'desc')
            ->where('state', '>', 0)
            ->select();


        $this->assign(['comment' => $comment, 'shuo' => $shuo]);
        return view('page/time_line');
    }

    /**
     * 评论添加
     * @param Request $request
     * @throws \Exception
     */
    public function commentAdd(Request $request)
    {
        $comment = $request->only(['nickname', 'pid' => 0, 'content', 'email', 'site' => null, 'article_id' => null, 'page_id' => null]);

        $comment['os'] = getOS();
        $comment['browser'] = getBrowser();

        $result = model('CommentModel')->add($comment);

        if (1 !== $result) {
            $this->error($result);
        }

        $title = !empty($comment['article_id'])
            ? model('ArticleModel')->find($comment['article_id'])['title']
            : config("page.{$comment['page_id']}")[0];

        $link = !empty($comment['article_id'])
            ? "https://muxik.top/info/{$comment['article_id']}.html"
            : config("page.{$comment['page_id']}")[0] . ".html";

        if (0 != $comment['pid']) {
            // 发送邮件 : 回复
            $content = model('CommentModel')->find($comment['pid']);
            $template = getReplyTemplate($content['nickname'], $content['content'], $comment['content'], $comment['nickname'], $link, $title);

            mailto($content['email'], $template, config('comment.reply_title'));
        }

        // 发送邮件 : 评论
        $template = getCommentTemplate($comment['nickname'], $comment['email'], $link, $title, $comment['content']);
        mailto(config('comment.email'), $template, config('comment.comment_title'));

        $this->success('评论成功!');
    }

}