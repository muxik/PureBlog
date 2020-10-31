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
     * @param $data
     * @return array|bool|string
     */
    public function add($data)
    {
        // 数据验证
        $validate = new \app\common\validate\Category();
        if (!$validate->scene('add')->check($data)){
            return $validate->getError();
        }

        $result = $this->save($data);
        if ($result) return true;
        else return '添加失败，请稍后再试';
    }
}
