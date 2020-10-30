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

Route::group('admin', function () {
    Route::rule('/', 'admin/login/index', 'get');
    Route::rule('login', 'admin/login/login', 'post');
});

Route::group('admin',function (){
    Route::rule('index', 'admin/index/index', 'get');

    // 管理员资源路由
    Route::resource('user','admin/admin');

})->middleware('CheckLogin');


Route::get('test', function (){
    try {
        return json(\think\Db::name("admin")->select());
    } catch (Exception $e) {
    }
});