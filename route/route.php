<?php




// 登录
Route::get('admin/login','admin/Login/index');
Route::post('admin/login.do', 'admin/Login/login');


Route::group('admin', function (){
    Route::get('index','admin/index/index');
    Route::get('welcome','admin/index/welcome');

    Route::get('category', 'admin/category/index');
    Route::get('category/list', 'admin/category/lst');

})->middleware(['CheckLogin']);

Route::get('/', 'index/index/index');