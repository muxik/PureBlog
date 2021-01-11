<?php


namespace app\admin\controller;


use app\common\controller\index\IndexController;
use app\common\model\WebModel;
use think\Request;

class Web extends IndexController
{

    protected $model;

    protected function initialize()
    {
        $this->model = new WebModel();
        parent::initialize();
    }

    public function index()
    {
        $info = $this->model->find();
        return view()->assign(['web'=>$info]);
    }


    public function update(Request $request)
    {
        $info = $request->only(['motto','id','logo','name','description','email','copyright','keyword']);

        $result = $this->model->edit($info);

        if (1 !== $result) {
            $this->error($result);
        }
        $this->success('修改成功！');
    }

    public function changeClose(Request $request)
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



}