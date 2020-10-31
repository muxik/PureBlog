<?php

namespace app\admin\controller;

use think\Controller;
use think\Request;

class Admin extends Controller
{
    /**
     * 显示资源列表
     *
     * @return \think\Response
     */
    public function index()
    {
        // 搜索实现
        if (input('username')) {
            $admins = model('Admin')
                ->searchNameAttr(input('username'), 'username')
                ->order('super', 'desc')
                ->paginate(5);
            return view()->assign(['admins' => $admins]);
        }

        $admins = model('Admin')
            ->order('super', 'desc')
            ->paginate(5);
        return view()->assign(['admins' => $admins]);

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
        if (!session('admin.super') == true) $this->error("你不是超级管理员");

        $data = [
            'username' => $request->param('username'),
            'password' => md5($request->param('password')),
            'status' => $request->param('status', 1),
            'super' => $request->param('super', 0),
        ];
        $result = model('Admin')->add($data);
        if ($result === true) $this->success('添加成功！');
        // TODO 功能测试 muxi_k
        $this->error($result);
    }

    /**
     * 显示指定的资源
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
        if (!session('admin.super') == true) $this->error("你不是超级管理员");

        $admin = model('Admin')
            ->find($id);
        return view()->assign(['admin' => $admin]);
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
            'password' => $request->param('password'),
            'status' => $request->param('status', 1),
            'super' => $request->param('super', 0),
        ];

        $result = model('Admin')->edit($id, $data);
        if ($result === true) $this->success('修改成功', '/admin/user');
        else $this->error($result);
    }

    /**
     * 删除指定资源
     *
     * @param int $id
     * @return \think\Response
     * @throws \Exception
     */
    public function delete($id)
    {
        if (!session('admin.super') == true) $this->error("你不是超级管理员");


        $result = model('Admin')->find($id)->delete();
        if ($result === true) $this->success('删除成功！', '/admin/user');
        else $this->error('删除失败，请重试！', '/admin/user');
    }
}
