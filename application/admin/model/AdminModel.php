<?php

namespace app\admin\model;

use app\admin\validate\AdminValidate;
use think\Model;
use think\model\concern\SoftDelete;

class AdminModel extends Model
{
    use SoftDelete;

    protected $name = "admin";

    public function login($data)
    {
        // 数据验证
        $validate = new AdminValidate();
        if (!$validate->scene('login')->check($data)){
            return $validate->getError();
        }

        // 密码加密
        $user = $this->where(['username' => $data['username']])->find();
        if (!$user) return "帐号或密码错误";



        $password = md5($data['password'] .  salt($user['id']));

        // 返回结果
        $result = $this->where(['username' => $user['username'], 'password' => $password])->find();

        // 记住密码
        if (!empty($data['rememberMe'])){
            cookie(['prefix' => 'admin_', 'expire' => 259200]);
            cookie('user', ['id' => $result['id'], 'username' => $result['username']], 259200);
        }


        // 保存session
        session('user', ['username' => $result['username'], 'id' => $result['id']], 'admin_');


        if (!$result) return "帐号或密码错误";
        return 1;
    }
}
