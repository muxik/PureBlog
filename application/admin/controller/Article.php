<?php

namespace app\admin\controller;

use think\Controller;
use think\Request;

class Article extends Controller
{
    /**
     * 显示资源列表
     *
     * @return \think\Response
     */
    public function index()
    {
        $article = model('Article')
            ->order('create_time', 'asc')
            ->with(['category', 'admin'])
            ->paginate(10);
        return view()->assign(['articles' => $article]);
    }

    /**
     * 显示创建资源表单页.
     *
     * @return \think\Response
     */
    public function create()
    {
        return  view();
    }

    /**
     * 保存新建的资源
     *
     * @param  \think\Request  $request
     * @return \think\Response
     */
    public function save(Request $request)
    {
        $data = [
            'title' => $request->param('title'),
            'top' => $request->param('title'),
            'u_id' => $request->param('title'),
            'tag' => $request->param('title'),
            'pic' => $request->param('title'),
            'content' => $request->param('title'),
            'state' => $request->param('title'),
        ];
    }

    /**
     * 显示指定的资源
     *
     * @param  int  $id
     * @return \think\Response
     */
    public function read($id)
    {
        //
    }

    /**
     * 显示编辑资源表单页.
     *
     * @param  int  $id
     * @return \think\Response
     */
    public function edit($id)
    {
        //
    }

    /**
     * 保存更新的资源
     *
     * @param  \think\Request  $request
     * @param  int  $id
     * @return \think\Response
     */
    public function update(Request $request, $id)
    {

    }

    /**
     * 删除指定资源
     *
     * @param  int  $id
     * @return \think\Response
     */
    public function delete($id)
    {
        //
    }

    /**
     * 文件上传
     * @param Request $request
     */
    public function upload(Request $request)
    {
        $file = $request->file('file');
        $info = $file->move('./static/upload');
        if ($info){
            $this->success('上传成功','/static/upload/'. $info->getSaveName());
        }
    }
}
