<?php


namespace app\admin\validate;


use think\Validate;

class LinkValidate extends Validate
{

    /**
     * 定义验证规则
     * 格式：'字段名'	=>	['规则1','规则2'...]
     *
     * @var array
     */
    protected $rule = [
        'link|链接' => 'require',
        'description|描述' => 'require',
        'logo|LOGO' => 'require',
        'title|' => 'require',
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
        return $this->only(['link', 'description', 'logo','title']);
    }
}