<?php

namespace app\admin\controller;

use app\admin\model\AdminModel;
use think\Request;

class Admin extends AdminController
{
    protected $model;

    protected function initialize()
    {
        $this->model = new AdminModel();
        parent::initialize();
    }

    public function index()
    {
        return view()->assign([
            'admin' => $this->model->paginate(10),
            'count' => $this->model->count()
        ]);
    }

    public function create()
    {
        return view();
    }

    public function edit($id)
    {
        $admin = $this->model->where('id', $id)->find();

        return view()->assign([
            'admin' => $admin
        ]);
    }

    public function add(Request $request)
    {
        $info = $request->only(['username', 'nickname', 'email', 'password']);
        $info['state'] = 1;
        $result = $this->model->add($info);

        if (1 !== $result) {
            $this->error($result);
        }
        $this->success('添加成功！');
    }

    public function update(Request $request)
    {
        $info = $request->only(['email', 'nickname', 'password', 'id']);
        $result = $this->model->edit($info);
        if (1 !== $result) {
            $this->error($result);
        }
        $this->success('修改成功！');
    }

    public function changeState()
    {
        $this->change($this->model,'state');
    }

    public function del()
    {
        $this->delete($this->model);
    }

}
