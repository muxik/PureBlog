<?php

namespace app\common\model;

use think\Model;

class Web extends Model
{
    protected $name = "web";

    public function edit($data)
    {
        $result = $this
            ->where('id', 1)
            ->update($data);
        if ($result) return true;
        else return "服务器错误，请稍候再试!";
    }

}
