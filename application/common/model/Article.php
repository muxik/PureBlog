<?php

namespace app\common\model;

use think\Model;
use think\model\concern\SoftDelete;

class Article extends Model
{
    // 软删除
    use SoftDelete;

    // 表名
    protected $name = "article";

    // 只读字段
    protected $readonly = ['read'];

    /**
     * 关联栏目
     * @return \think\model\relation\BelongsTo
     */
    public function category(){
        return $this->belongsTo('Category', 'category_id', 'id');
    }

    /**
     * 关联作者
     * @return \think\model\relation\BelongsTo
     */
    public function admin(){
        return $this->belongsTo('Admin', 'u_id', 'id');
    }


}
