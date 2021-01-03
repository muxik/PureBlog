<?php


namespace app\common\validate;


use think\Validate;

class ArticleValidate extends Validate
{

    /**
     * 定义验证规则
     * 格式：'字段名'	=>	['规则1','规则2'...]
     *
     * @var array
     */
    protected $rule = [
        'title|标题' => 'require',
        'description|描述' => 'require',
        'content|内容' => 'require',
        'pic|封面图' => 'require',
        'tag|标签' => 'require',
        'category_id|分类' => 'require'
    ];

    /**
     * 定义错误信息
     * 格式：'字段名.规则名'	=>	'错误信息'
     *
     * @var array
     */
    protected $message = [];


    /** 添加文章验证场景
     * @return ArticleValidate
     */
    public function sceneAdd()
    {
        return $this->only(['title', 'description', 'content', 'pic', 'tag', 'require']);
    }

}