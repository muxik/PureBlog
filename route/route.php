<?php




// 登录
Route::get('admin/login','admin/Login/index');
Route::post('admin/login.do', 'admin/Login/login');


Route::group('admin', function (){
    Route::get('index','admin/index/index');

})->middleware(['CheckLogin']);

Route::get('/', 'index/index/index');