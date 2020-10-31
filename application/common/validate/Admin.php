<?php

namespace app\common\validate;

use think\Validate;

class Admin extends Validate
{
    /**
     * 定义验证规则
     * 格式：'字段名'    =>    ['规则1','规则2'...]
     *
     * @var array
     */
    protected $rule = [
        'username|管理员账户' => 'require',
        'password|管理员密码' => 'require',
        'code|验证码' => 'require|captcha'
    ];


    /**
     * 登录验证
     * @return Admin
     */
    public function sceneLogin()
    {
        return $this->only(['username', 'password', 'code']);
    }

    /**
     * 添加管理员验证
     * @return Admin
     */
    public function sceneAdd()
    {
        return $this->only(['username','password'])
            ->append(['username|用户名' => 'unique:admin']);
    }

    /**
     * 管理员修改验证
     * @return Admin
     */
    public function sceneEdit()
    {
        return $this->only(['password']);
    }
}
