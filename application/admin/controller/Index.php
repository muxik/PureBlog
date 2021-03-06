<?php

namespace app\admin\controller;

use app\common\controller\admin\AdminController;
use think\Controller;
use think\Db;

class Index extends AdminController
{

    protected function initialize()
    {

        parent::initialize();
    }

    public function index()
    {
        // 查询用户信息
        session(['prefix' => 'admin_',]);
        $id = session('user.id');

        if (empty($id) & !empty(cookie('admin_user')))
            $id = cookie('admin_user')['id'];

        $admin = model('AdminModel')
            ->where('id', $id)
            ->find();

        $head = md5(trim(strtolower($admin['email'])));

        return view()
            ->assign([
                'admin' => $admin,
                'head' => $head
            ]);
    }

    public function welcome()
    {

        $bytes = disk_free_space(".");
        $si_prefix = array( 'B', 'KB', 'MB', 'GB', 'TB', 'EB', 'ZB', 'YB' );
        $base = 1024;
        $class = min((int)log($bytes , $base) , count($si_prefix) - 1);

        $DISK_FREE_SPACE =  sprintf('%1.2f' , $bytes / pow($base,$class)) . $si_prefix[$class];
        $PURE_BLOG_VERSION = '2.0.0';
        $PHP_VERSION = PHP_VERSION;
        $MYSQL_VERSION =  Db::query("select VERSION()")[0]['VERSION()'];



        $viewData  = [
            'article_count' => model('ArticleModel')->count(),
            'user_count' => model('UserModel')->count(),
            'category_count' => model('CategoryModel')->count(),
            'link_count' => model('LinkModel')->count(),
            'admin_count' => model('AdminModel')->count(),
            'pure_blog_version' => $PURE_BLOG_VERSION,
            'php_version' => $PHP_VERSION,
            'mysql_version' => $MYSQL_VERSION,
            'disk_free_space' => $DISK_FREE_SPACE,
        ];

        return view()->assign($viewData);
    }
}
