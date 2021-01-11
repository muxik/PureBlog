<?php


namespace app\admin\controller;


use app\common\controller\admin\AdminController;
use app\common\model\ShuoModel;
use think\Request;

class Shuo extends AdminController
{

    protected $model;

    protected function initialize()
    {
        $this->model = new ShuoModel();
        parent::initialize();
    }

    public function index()
    {
        $count = $this->model
            ->count();
        $link = $this->model
            ->paginate(10);

        return view()->assign([
            'link' => $link,
            'count' => $count
        ]);
    }

    public function create()
    {
        return view();
    }

    public function add(Request $request)
    {
        $info = $request->only(['content', 'year' => date("Y")]);

        $result = $this->model->add($info);

        if (1 !== $result) {
            $this->error($result);
        }
        $this->success('添加成功！');
    }

    public function edit($id)
    {
        $info = $this->model->where('id', $id)->find();
        return view()->assign(['shuo' => $info]);
    }

    public function update(Request $request)
    {
        $info = $request->only(['content','id']);

        $result = $this->model->edit($info);

        if (1 !== $result) {
            $this->error($result);
        }
        $this->success('修改成功！');
    }

    public function changeState()
    {
        $this->change($this->model, 'state');
    }


    public function del()
    {
        $this->delete($this->model);
    }


}