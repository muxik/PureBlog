<?php

namespace app\common\validate;


use think\Validate;

class CategoryValidate extends Validate
{
    /**
     * 定义验证规则
     * 格式：'字段名'	=>	['规则1','规则2'...]
     *
     * @var array
     */	
	protected $rule = [
	    'name|栏目名称' => 'require',
        'pid|栏目分类' => 'require'
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
        return $this->only(['name', 'pid']);
    }

    public function sceneEdit()
    {
        return $this->only(['name']);
    }

}
