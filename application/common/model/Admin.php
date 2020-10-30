<?php

namespace app\common\model;

use think\Model;
use think\model\concern\SoftDelete;

class Admin extends Model
{
    use SoftDelete;

    protected $name = "admin";

    /**
     * 搜索
     * @param $keyword
     * @param $value
     */
    public function searchNameAttr($keyword, $value)
    {
        $where[] = [
            $value, 'like', '%' . $keyword . '%'
        ];
        return $this->where($where);
    }
}
