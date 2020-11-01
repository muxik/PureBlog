<?php

namespace app\index\controller;

class Index
{
    public function index()
    {
        $articles = model('Article')->select();
        return view()->assign(['articles' => $articles]);
    }
}
