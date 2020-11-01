<?php

namespace app\common\validate;

use think\Validate;

class Article extends Validate
{
    /**
     * 定义验证规则
     * 格式：'字段名'    =>    ['规则1','规则2'...]
     *
     * @var array
     */
    protected $rule = [
        'title|标题' => 'require',
        'content|内容' => 'require',
        'category_id|分类' => 'require',
        'pic|图片路径' => 'require',
        'tag|标签' => 'require'
    ];

    /**
     * 添加验证场景
     * @return Article
     */
    public function sceneAdd()
    {
        return $this->only(['title', 'content', 'category_id', 'pic', 'tag']);
    }

    /**
     * 更新验证场景
     * @return Article
     */
    public function sceneEdit()
    {
        return $this->only(['title', 'content', 'category_id', 'pic', 'tag']);
    }
}
