<?php

namespace app\admin\validate;

use think\Validate;

class AdminValidate extends Validate
{
    /**
     * 定义验证规则
     * 格式：'字段名'    =>    ['规则1','规则2'...]
     *
     * @var array
     */
    protected $rule = [
        'username|用户名' => 'require',
        'password|密码' => 'require',
        'captcha|验证码' => 'require|captcha'
    ];

    /**
     * 定义错误信息
     * 格式：'字段名.规则名'    =>    '错误信息'
     *
     * @var array
     */
    protected $message = [];


    /**
     * 登录场景
     * @return AdminValidate
     */
    public function sceneLogin()
    {
        return $this->only(['username', 'password', 'captcha']);
    }

    public function sceneAdd()
    {
        return $this->only(['username', 'password', 'nickname'])
            ->remove('captcha')
            ->append(['nickname|昵称' => 'require']);
    }


    public function sceneEdit()
    {
        return $this->only(['password'])->append(['nickname' => 'require', 'email' => 'require']);
    }
}
