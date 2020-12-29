<?php


namespace app\admin\controller;


use think\Controller;
use think\Request;

class Category extends Controller
{
    public function index()
    {
        $category = model('CategoryModel');
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


    public function add(Request $request)
    {

        $data = $request->only(['pid', 'name', 'sort' => 0]);
        $result = model('CategoryModel')->add($data);
        if (1 !== $result) {
             $this->error($result);
        }

        $this->success("添加成功！");
    }

    public function changeState(Request $request)
    {
        $state = $request->post('value', 1);
        $result = model('CategoryModel')->where('id', $request->post('id'))->update(['state' => $state]);
        if (!$result)
            $this->error("修改失败");
        else
            $this->success("修改成功！");
    }


    public function changeSort(Request $request)
    {
        $id  = $request->post('id');
        $sort = $request->post('sort',0);
        $result = model('CategoryModel')->where('id', $id)->update(['sort' => $sort]);
        if (!$result)
            $this->error("修改失败");
        else
            $this->success("修改成功！");
    }
}