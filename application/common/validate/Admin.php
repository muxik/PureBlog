<?php

namespace app\common\validate;

use think\Validate;

class Admin extends Validate
{
    /**
     * 定义验证规则
     * 格式：'字段名'	=>	['规则1','规则2'...]
     *
     * @var array
     */	
	protected $rule = [
	    'username|管理员账户' => 'require',
        'password|管理员密码' => 'require',
        'code|验证码'        => 'require|captcha'
    ];


    /**
     * 登录
     * @return Admin
     */
    public function sceneLogin()
    {
        return $this->only(['username', 'password', 'code']);
    }
}
