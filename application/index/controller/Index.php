<?php

namespace app\index\controller;

class Index extends IndexController
{
    public function index()
    {

        $viewData = [
            'article' => $this->article->paginate(5),
            'page' => $this->article->paginate(5)->toArray(),
            'new_article' => $this->article->limit(5),
            'category' => $this->category->select(),
            'web' => $this->web,
            'tag' => $this->getTags()
        ];

        // 搜索
        if (!empty(input('q'))) {
            $viewData['article'] = $this->article
                ->where([['title', 'like', '%' . input('q') . '%']])
                ->paginate(5);
            return view()->assign($viewData);
        }

        return view()->assign($viewData);
    }

    protected function getTags()
    {
        $tag = [];
        $this->tags = array_unique($tag);
        foreach (
            model('ArticleModel')
                ->field('tag')
                ->where('state', '>', 0)
                ->select()
            as
            $key => $vo
        ) {
            foreach (explode(',', $vo['tag']) as $k => $v) {
                array_push($tag, $v);
            }
        }
        return $tag;
    }

    public function category()
    {
        $viewData = [
            'article' => $this->article->paginate(5),
            'page' => $this->article->paginate(5)->toArray(),
            'new_article' => $this->article->limit(5),
            'category' => $this->category->select(),
            'web' => $this->web,
            'tag' => $this->getTags()
        ];

        $category_name = '';
        foreach (explode('/', request()->url()) as $key => $vo) {
            if ($key > 1) {
                $line = $key > 2 ? '/' : '';
                $category_name .= $line . $vo;
            }
        }

        $category = $this->category->where('page', $category_name)->find();
        $id = 0;
        if (!empty($category['id'])){
            $id = $category['id'];
        }
        $viewData['article'] = $this->article
            ->where('category_id', $id)
            ->paginate(5);
        return view('index')->assign($viewData);
    }

    public function like()
    {
        $this->web->setInc('like');
    }
}
