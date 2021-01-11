<?php

namespace app\index\controller;

use app\common\controller\index\IndexController;

class Index extends IndexController
{
    public function index()
    {
        $this->assign([
            'page' => $this->article->paginate(5)->toArray()
        ]);

        // 搜索
        if (!empty(input('q'))) {
            $viewData['article'] = $this->article
                ->where([['title', 'like', '%' . input('q') . '%']])
                ->paginate(5);
            return view()->assign($viewData);
        }

        return view();
    }

    public function category()
    {
        $this->assign([
            'page' => $this->article->paginate(5)->toArray()
        ]);


        $category_name = '';
        foreach (explode('/', request()->url()) as $key => $vo) {
            if ($key > 1) {
                $line = $key > 2 ? '/' : '';
                $category_name .= $line . $vo;
            }
        }

        $category = $this->category->where('page', $category_name)->find();
        $id = 0;
        if (!empty($category['id'])) {
            $id = $category['id'];
        }
        $viewData['article'] = $this->article
            ->where('category_id', $id)
            ->paginate(5);

        return view('index');
    }



    public function like()
    {
        $this->web->setInc('like');
    }
}
