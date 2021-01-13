<?php


namespace app\admin\controller;


use app\common\controller\admin\AdminController;
use app\common\model\CommentModel;

class Comment extends AdminController
{

    protected $model;

    protected function initialize()
    {
        $this->model = new CommentModel();
        parent::initialize();
    }

    public function index()
    {
        $this->assign([
            'count' => $this->model->count(),
            'comment' => $this->model->paginate(10)
            ]);
        return view();
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