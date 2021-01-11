<?php


namespace app\index\controller;


use app\common\controller\index\IndexController;
use think\Request;

class Page extends IndexController
{

    /**
     * 关于
     * @return \think\response\View
     */
    public function about()
    {
        $comment = model('CommentModel')->where('page_id', 1)->select();
        $this->assign(['comment' => $comment]);
        return view('page/about');
    }

    /**
     * 友人
     * @return \think\response\View
     */
    public function link()
    {
        $comment = model('CommentModel')->where('page_id', 2)->select();
        $this->assign(['comment' => $comment]);
        return view('page/link');
    }

    /**
     * 碎碎念
     * @return \think\response\View
     */
    public function timeLine()
    {
        $comment = model('CommentModel')->where('page_id', 3)->select();
        $shuo = model('ShuoModel')->order('create_time','desc')->where('state', '>', 0)->select();


        $this->assign(['comment' => $comment, 'shuo' => $shuo]);
        return view('page/time_line');
    }

    /**
     * 评论添加
     * @param Request $request
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

        $this->success('评论成功!');
    }

}