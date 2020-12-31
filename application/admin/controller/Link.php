<?php


namespace app\admin\controller;


use app\admin\model\LinkModel;
use think\Controller;
use think\Request;

class Link extends Controller
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

    public function edit()
    {
        return view();
    }

    public function update(Request $request)
    {

    }

    public function changeState(Request $request)
    {
        $state = $request->post('value', 1);
        $result = $this
            ->model
            ->where('id', $request->post('id'))
            ->update(['state' => $state]);
        if (!$result)
            $this->error("修改失败");
        else
            $this->success("修改成功！");

    }


    public function del()
    {

    }


}