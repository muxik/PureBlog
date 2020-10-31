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
        $categorys = model('Category')->order('sort', 'asc')->select();
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
        if (model('Category')->find($id)->delete())
            $this->success('删除成功！');
        else $this->error('删除失败，请稍后再试！');
    }
}
