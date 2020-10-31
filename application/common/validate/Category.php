<?php

namespace app\common\validate;

use think\Validate;

class Category extends Validate
{
    /**
     * 定义验证规则
     * 格式：'字段名'	=>	['规则1','规则2'...]
     *
     * @var array
     */	
	protected $rule = [
	    'name|名称' => 'require',
        'sort|排序' => 'require|number'
    ];

    /**
     * 添加栏目验证场景
     * @return Category
     */
	public function sceneAdd()
    {
        return $this->only(['name','sort']);
    }
}
