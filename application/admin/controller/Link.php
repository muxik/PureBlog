<?php


namespace app\admin\controller;


use app\common\model\LinkModel;
use think\Request;

class Link extends AdminController
{

    protected $model;

    protected function initialize()
    {
        $this->model = new LinkModel();
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
        $info = $request->only(['title','logo', 'link', 'description']);

        $result = $this->model->add($info);

        if (1 !== $result) {
            $this->error($result);
        }
        $this->success('添加成功！');
    }

    public function edit($id)
    {
        $info = $this->model->where('id', $id)->find();
        return view()->assign(['link' => $info]);
    }

    public function update(Request $request)
    {
        $info = $request->only(['id','title','logo', 'link', 'description']);

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

    public function setConnect()
    {
        $this->change($this->model, 'type');
    }


    public function del()
    {
        $this->delete($this->model);
    }


}