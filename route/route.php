<?php


// 登录
Route::get('admin/login', 'admin/Login/index');
Route::post('admin/login.do', 'admin/Login/login');


Route::group('admin', function () {
    Route::get('index', 'admin/index/index');
    Route::get('welcome', 'admin/index/welcome');
    Route::get('out', function (){
        session(null);
        cookie(null,'admin_');
        return redirect('/admin/login');
    });

    // 分类模块
    Route::group('category', function () {
        Route::get('/', 'admin/category/index');
        Route::get('edit/:id', 'admin/category/edit');
        Route::post('update', 'admin/category/update');
        Route::post('add', 'admin/category/add');
        Route::post('changeSort', 'admin/category/changeSort');
        Route::post('changeState', 'admin/category/changeState');
        Route::post('setPage', 'admin/category/setPage');
        Route::post('del', 'admin/category/del');
    });

    //文章模块
    Route::group('article', function () {
        Route::get('/', 'admin/article/index');
        Route::get('create', 'admin/article/create');
        Route::get('edit/:id', 'admin/article/edit');
        Route::post('add', 'admin/article/add');
        Route::post('update', 'admin/article/update');
        Route::post('del', 'admin/article/del');
        Route::post('changeState', 'admin/article/changeState');
        Route::post('changeTop', 'admin/article/changeTop');
        Route::post('upload', 'admin/article/upload');
    });

    // 友情链接模块
    Route::group('link',function (){
        Route::get('/', 'admin/link/index');
        Route::get('create', 'admin/link/create');
        Route::get('edit/:id', 'admin/link/edit');
        Route::post('del', 'admin/link/del');
        Route::post('add','admin/link/add');
        Route::post('update','admin/link/update');
        Route::post('changeState', 'admin/link/changeState');
        Route::post('setConnect', 'admin/link/setConnect');
        Route::post('del', 'admin/link/del');
    });

    // 管理员模块

    Route::get('admin/edit/:id', 'admin/admin/edit');
    Route::get('admin/create', 'admin/admin/create');
    Route::get('admin', 'admin/admin/index');
    Route::group('admin',function (){
        Route::post('del', 'admin/admin/del');
        Route::post('add','admin/admin/add');
        Route::post('update','admin/admin/update');
        Route::post('changeState', 'admin/admin/changeState');
        Route::post('del', 'admin/admin/del');
    })->middleware(['Auth']);

    // 用户模块
    Route::group('user',function (){
        Route::get('edit/:id', 'admin/user/edit');
        Route::get('create', 'admin/user/create');
        Route::get('/', 'admin/user/index');
        Route::post('del', 'admin/user/del');
        Route::post('add','admin/user/add');
        Route::post('update','admin/user/update');
        Route::post('changeState', 'admin/user/changeState');
    });

    // 网站信息模块
    Route::group('web',function (){
        Route::get('/', 'admin/web/index');
        Route::post('update','admin/web/update');
        Route::post('changeClose', 'admin/web/changeClose');
    });


})->middleware(['CheckLogin']);

Route::get('/', 'index/index/index');

Route::get('info/:id','index/index/info')->ext('html');

Route::get('category/[:page]','index/index/category');
Route::get('about', 'index/page/about')->ext('html');
Route::get('link', 'index/page/link')->ext('html');

Route::post('muxik/like', 'index/index/like');
Route::post('commentAdd', 'index/page/commentAdd');