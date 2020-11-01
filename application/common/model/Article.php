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
    public function category()
    {
        return $this->belongsTo('Category', 'category_id', 'id');
    }

    /**
     * 关联作者
     * @return \think\model\relation\BelongsTo
     */
    public function admin()
    {
        return $this->belongsTo('Admin', 'u_id', 'id');
    }


    /**
     * 添加 模型
     * @param $data
     * @return array|bool|string
     */
    public function add($data)
    {

        // 验证数据
        $validate = new \app\common\validate\Article();
        if (!$validate->scene('add')->check($data)) {
            return $validate->getError();
        }

        $result = $this->save($data);
        if ($result) return true;
        else return "服务器错误请,稍后再试!";
    }


    public function edit($id, $data)
    {
        $validate = new \app\common\validate\Article();
        if (!$validate->scene('edit')->check($data)) {
            return $validate->getError();
        }

        $result = $this->where('id', $id)->update($data);
        if ($result) return true;
        else return "服务器错误请稍后再试";
    }

}
