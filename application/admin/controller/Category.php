<?php

namespace app\admin\controller;

use think\Controller;
use think\Request;

class Category extends Controller
{
    /**
     * 显示资源列表
     *
     * @return \think\Response
     */
    public function index()
    {
        $categorys = model('Category')
            ->order('sort', 'asc')
            ->with('article')
            ->paginate(5);
        return view()->assign(['categorys' => $categorys]);
    }

    /**
     * 显示创建资源表单页.
     *
     * @return \think\Response
     */
    public function create()
    {
        return view();
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
            'name' => $request->param('name'),
            'sort' => $request->param('sort', 1),
            'state' => $request->param('state', 0)
        ];

        $result = model('Category')->add($data);
        if ($result === true) $this->success('添加成功！');
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
        //
    }

    /**
     * 显示编辑资源表单页.
     *
     * @param int $id
     * @return \think\Response
     */
    public function edit($id)
    {
        $category = model('Category')->find($id);
        return view()->assign(['cate' => $category]);
    }

    /**
     * 栏目 显示/隐藏
     * @param Request $request
     * @param $id
     */
    public function updateState(Request $request, $id)
    {
        $result = model('Category')->updateState($id, $request->param('status'));
        if ($result === true) $this->success('状态更新完成', '/admin/user');
        else $this->error($result);
    }

    /**
     * 栏目排序
     * @param Request $request
     * @param $id
     */
    public function updateSort(Request $request, $id)
    {
        $result = model('Category')->updateSort($id, $request->param('sort'));
        if ($result === true) $this->success('排序成功', '/admin/cate');
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
        $data = [
            'state' => $request->param('state'),
            'name' => $request->param('name'),
            'sort' => $request->param('sort'),
        ];

        $result = model('Category')->edit($id, $data);

        if ($result === true) $this->success('添加成功');
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
        $cate = model('Category')
            ->with('article')
            ->find($id);

        $result = $cate->together('article')->delete();
        if ($result)
            $this->success('删除成功！');
        else
            $this->error('删除失败，请稍后再试！');
    }
}
