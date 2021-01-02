<?php


namespace app\admin\validate;


use think\Validate;

class WebValidate extends Validate
{

    /**
     * 定义验证规则
     * 格式：'字段名'	=>	['规则1','规则2'...]
     *
     * @var array
     */
    protected $rule = [
        'name|网站名称' => 'require',
        'description|网站描述' => 'require',
        'logo|网站LOGO' => 'require',
        'copyright|版权信息' => 'require',
    ];

    /**
     * 定义错误信息
     * 格式：'字段名.规则名'	=>	'错误信息'
     *
     * @var array
     */
    protected $message = [];

    public function sceneEdit()
    {
        return $this->only(['link', 'description', 'logo','title']);
    }
}