<?php




// 登录
Route::get('admin/login','admin/Login/index');
Route::post('admin/login.do', 'admin/Login/login');


Route::group('admin', function (){
    Route::get('index','admin/index/index');
    Route::get('welcome','admin/index/welcome');

    Route::get('category', 'admin/category/index');
    Route::post('category/add', 'admin/category/add');
    Route::post('category/changeSort','admin/category/changeSort');
    Route::post('category/changeState','admin/category/changeState');

})->middleware(['CheckLogin']);

Route::get('/', 'index/index/index');