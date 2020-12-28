<?php
namespace app\index\controller;

class Index
{
    public function index()
    {
        return md5('admin' . salt(1));
    }
}
