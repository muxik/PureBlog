<?php

namespace app\common\exception;

use Exception;
use think\exception\Handle;
use think\exception\HttpException;

class Http extends Handle
{

    public function render(Exception $e)
    {
        if ($e instanceof HttpException) {
            if (stristr($e->getMessage(), "module not exists:")) {
                return view('error');
            }

            if (strstr($e->getMessage(), 'controller not exists:')){
                return view('../../view/error');
            }
        }

        //可以在此交由系统处理
        return parent::render($e);
    }

}