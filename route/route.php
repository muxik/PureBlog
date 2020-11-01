<?php
// +----------------------------------------------------------------------
// | ThinkPHP [ WE CAN DO IT JUST THINK ]
// +----------------------------------------------------------------------
// | Copyright (c) 2006~2018 http://thinkphp.cn All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( http://www.apache.org/licenses/LICENSE-2.0 )
// +----------------------------------------------------------------------
// | Author: liu21st <liu21st@gmail.com>
// +----------------------------------------------------------------------

Route::group('/',function (){
    Route::rule('', 'index/index/index', 'get');
});



Route::group('admin', function () {
    Route::rule('/', 'admin/login/index', 'get');
    Route::rule('login', 'admin/login/login', 'post');
});

Route::group('admin',function (){
    Route::rule('index', 'admin/index/index', 'get');
    Route::rule('out', 'admin/index/out', 'get');

    // 管理员管理 资源路由
    Route::resource('user','admin/admin');
    Route::post('user/up_state/:id','admin/admin/updateState');

    // 栏目管理 资源路由
    Route::resource('cate','admin/category');
    Route::post('cate/up_state/:id','admin/category/updateState');
    Route::post('cate/up_sort/:id','admin/category/updateSort');

    // 文章管理 资源路由
    Route::resource('article','admin/article');
    Route::post('upload','admin/article/upload');
    Route::post('article/up_state/:id','admin/article/updateState');
    Route::post('article/up_top/:id','admin/article/updateTop');

    // 网站设置
    Route::get('web','admin/web/index');
    Route::post('update','admin/web/update');



})->middleware('CheckLogin');


Route::get('test', function (){
    try {
        return json(\think\Db::name("category")->select());
    } catch (Exception $e) {
    }
});