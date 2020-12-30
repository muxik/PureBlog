<?php


// 登录
Route::get('admin/login', 'admin/Login/index');
Route::post('admin/login.do', 'admin/Login/login');


Route::group('admin', function () {
    Route::get('index', 'admin/index/index');
    Route::get('welcome', 'admin/index/welcome');

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
        Route::get('edit', 'admin/article/edit');
        Route::post('add', 'admin/article/add');
        Route::post('update', 'admin/article/update');
        Route::post('changeTop', 'admin/article/changeTop');
        Route::post('changeState', 'admin/article/changeState');
        Route::post('upload', 'admin/article/upload');
    });



})->middleware(['CheckLogin']);

Route::get('/', 'index/index/index');
Route::get('/md5', function (){
    return md5('admin');
});