<?php

namespace app\common\model;

use think\Model;
use think\model\concern\SoftDelete;

/**
 * Class Admin
 * @package app\common\model
 * @author muxi_k
 */
class Admin extends Model
{
    use SoftDelete;

    protected $name = "admin";

    /**
     * 添加管理员
     * @param $data
     * @return array|string|boolean
     */
    public function add($data)
    {
        $validate = new \app\common\validate\Admin();
        if (!$validate->scene('add')->check($data)) {
            return $validate->getError();
        }

        $result = $this->save($data);
        if ($result) {
            return true;
        } else {
            return "添加失败";
        }
    }

    /**
     * 管理员修改
     * @param $id
     * @param $data
     * @return array|string|boolean
     */
    public function edit($id, $data)
    {
        // 验证数据
        $validate = new \app\common\validate\Admin();
        if (!$validate->scene('edit')->check($data)) {
            return $validate->getError();
        }

        // 更新数据
        $admin = $this->find($id);
        $admin->password = $data['password'];
        $admin->status = $data['status'];
        $admin->super = $data['super'];
        $result = $admin->save();
        if ($result) {
            return true;
        } else {
            return "修改失败";
        }
    }

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
