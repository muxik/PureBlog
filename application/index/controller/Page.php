<?php


namespace app\index\controller;


use think\Request;

class Page extends IndexController
{

    public function about()
    {
        $comment = model('CommentModel')->where('page_id', 1)->select();
        $this->assign(['comment' => $comment]);
        return view('page/about');
    }

    public function link()
    {
        $comment = model('CommentModel')->where('page_id', 2)->select();
        $this->assign(['comment' => $comment]);
        return view('page/link');
    }

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