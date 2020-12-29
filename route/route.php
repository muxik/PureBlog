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
        Route::post('del', 'admin/category/del');
    });


    //


})->middleware(['CheckLogin']);

Route::get('/', 'index/index/index');