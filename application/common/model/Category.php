<?php

namespace app\common\model;

use think\Model;
use think\model\concern\SoftDelete;

/**
 * Class Category
 * @package app\common\model
 * @author muxi_k
 */
class Category extends Model
{
    // 软删除
    use SoftDelete;

    // 表名
    protected $name = "category";

    /**
     * 添加栏目
     */
}
