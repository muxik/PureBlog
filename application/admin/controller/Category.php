<?php


namespace app\admin\controller;


use think\Controller;

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
            ->assign(['category' => $result]);
    }
}