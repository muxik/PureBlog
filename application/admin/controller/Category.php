<?php


namespace app\admin\controller;


use app\common\controller\admin\AdminController;
use app\common\model\CategoryModel;
use think\Request;

class Category extends AdminController
{

    protected $model;

    protected function initialize()
    {
        $this->model = new CategoryModel();
        parent::initialize();
    }

    public function index()
    {
        $category = $this->model;
        $tree = $category->order('sort', 'asc')->select();

        // 分类排序
        $result = [];
        for ($i = 0; $i < $category->count(); $i++) {
            if ($tree[$i]['pid'] == 0) {
                array_push($result, $tree[$i]);
                for ($k = 0; $k < $category->count(); $k++) {
                    if ($tree[$k]['pid'] == $tree[$i]['id']) array_push($result, $tree[$k]);
                }
            }
        }

        return view('category/index')
            ->assign([
                'category' => $result,
                'count' => $tree->count(),
                'p_category' => $category->where('pid', 0)->select()
            ]);
    }


    /**
     * 添加分类
     * @param Request $request
     */
    public function add(Request $request)
    {

        $data = $request->only(['pid', 'page', 'name', 'sort' => 0]);
        $result = $this->model->add($data);
        if (1 !== $result) {
            $this->error($result);
        }

        $this->success("添加成功！");
    }

    /**
     * @return string
     */
    public function edit($id)
    {
        return view()
            ->assign([
                'category' => $this->model->where('id', $id)->find(),
                'p_category' => $this->model->where('pid', 0)->select()
            ]);
    }


    public function update(Request $request)
    {
        if (empty($request->post('link'))) {
            $param = ['id', 'name', 'page' => 0];
        } else {
            $param = ['id', 'name', 'link' => null];
        }


        $data = $request->only($param);

        $result = $this->model->edit($data);
        if (1 !== $result) {
            $this->error($result);
        }
        $this->success('修改成功！');
    }

    /**
     * 改变分类状态
     */
    public function changeState()
    {
        $this->change($this->model, 'state');
    }


    /**
     *  改变分类排除
     */
    public function changeSort()
    {
        $this->change($this->model, 'sort', 'sort');
    }


    /**
     * 设置为单页
     */
    public function setPage()
    {
        $this->change($this->model, 'type');
    }


    /**
     * 删除分类
     * @param Request $request
     * @throws \Exception
     */
    public function del(Request $request)
    {
        $id = $request->post('id');

        $cateInfo = $this->model->with('article,article.comments')->find($id);
        foreach ($cateInfo['article'] as $k => $v) {
            $v->together('comments')->delete();
        }
        $result = $cateInfo->together('article')->delete();

        if (!$result)
            $this->error("删除失败");
        else
            $this->success("删除成功！");
    }
}