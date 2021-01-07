<?php


namespace app\common\validate;


use think\Validate;

class CommentValidate extends Validate
{
    /**
     * 定义验证规则
     * 格式：'字段名'	=>	['规则1','规则2'...]
     *
     * @var array
     */
    protected $rule = [
        'email|邮箱' => 'require',
        'content|内容' => 'require',
        'nickname|昵称' => 'require'
    ];

    /**
     * 定义错误信息
     * 格式：'字段名.规则名'	=>	'错误信息'
     *
     * @var array
     */
    protected $message = [];


    public function sceneAdd()
    {
        return $this->only(['email', 'content',' nickname']);
    }
}