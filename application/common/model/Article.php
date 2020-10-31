<?php

namespace app\common\model;

use think\Model;
use think\model\concern\SoftDelete;

class Article extends Model
{
    // 软删除
    use SoftDelete;

    // 表名
    protected $name = "article";

    // 只读字段
    protected $readonly = ['read'];


}
