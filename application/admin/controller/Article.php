<?php

namespace app\admin\controller;

use think\Controller;
use think\Request;

class Article extends Controller
{
    /**
     * 显示资源列表
     *
     * @return \think\Response
     */
    public function index()
    {
        $article = model('Article')
            ->order('create_time', 'asc')
            ->with(['category', 'admin'])
            ->paginate(10);
        return view()->assign(['articles' => $article]);
    }

    /**
     * 显示创建资源表单页.
     *
     * @return \think\Response
     */
    public function create()
    {
        $categorys = model('Category')->select();
        return view()->assign(
            ['categorys' => $categorys]
        );
    }

    /**
     * 保存新建的资源
     *
     * @param \think\Request $request
     * @return \think\Response
     */
    public function save(Request $request)
    {
        $data = [
            'title' => $request->param('title'),
            'desc' => $request->param('desc'),
            'top' => $request->param('top', 0),
            'u_id' => session('admin.id'),
            'tag' => $request->param('tag'),
            'pic' => $request->param('pic'),
            'content' => $request->param('content'),
            'state' => $request->param('state', 1),
            'category_id' => $request->param('category_id'),
        ];

        $result = model('Article')->add($data);
        if ($result === true) $this->success('文章添加成功!', '/admin/article');
        else $this->error($result);

    }

    /**
     * 显示指定的资源
     *
     * @param int $id
     * @return \think\Response
     */
    public function read($id)
    {
    }

    /**
     * 显示编辑资源表单页.
     *
     * @param int $id
     * @return \think\Response
     */
    public function edit($id)
    {
        $article = model('Article')->find($id);
        $categorys = model('Category')->select();

        return view()->assign([
            'article' => $article,
            'categorys' => $categorys
        ]);
    }

    /**
     * 文章状态更新
     * @param Request $request
     * @param $id
     */
    public function updateState(Request $request, $id){
        $result = model('Article')->updateState($id, $request->param('status'));
        if ($result === true) $this->success('状态更新完成', '/admin/user');
        else $this->error($result);
    }

    /**
     * 文章置顶更新
     * @param Request $request
     * @param $id
     */
    public function updateTop(Request $request, $id){
        $result = model('Article')->updateTop($id, $request->param('status'));
        if ($result === true) $this->success('操作成功！', '/admin/user');
        else $this->error($result);
    }

    /**
     * 保存更新的资源
     *
     * @param \think\Request $request
     * @param int $id
     * @return \think\Response
     */
    public function update(Request $request, $id)
    {
        $data = $request->except(['_method', 'file']);
        $data['state'] = $request->param('state', 0);
        $data['top'] = $request->param('top', 0);
        $result = model('Article')->edit($id, $data);

        if ($result) $this->success('修改完成');
        else $this->error($result);
    }

    /**
     * 删除指定资源
     *
     * @param int $id
     * @return \think\Response
     */
    public function delete($id)
    {
        $result = model('Article')->find($id)->delete();
        if ($result === true) $this->success('删除成功！', '/admin/article');
        else $this->error('删除失败，请重试！', '/admin/article');
    }

    /**
     * 文件上传
     * @param Request $request
     */
    public function upload(Request $request)
    {
        $file = $request->file('file');
        $info = $file->move('./static/upload');
        if ($info) {
            $this->success('上传成功', '/static/upload/' . $info->getSaveName());
        }
    }
}
