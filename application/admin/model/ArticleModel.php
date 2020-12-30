<?php


namespace app\admin\model;


use app\admin\validate\ArticleValidate;
use think\Model;
use think\model\concern\SoftDelete;

class ArticleModel extends Model
{
    use SoftDelete;

    protected $name = "article";


    /**
     * 关联栏目
     * @return \think\model\relation\BelongsTo
     */
    public function category()
    {
        return $this->belongsTo('CategoryModel', 'category_id', 'id');
    }

    /**
     * 关联作者
     * @return \think\model\relation\BelongsTo
     */
    public function admin()
    {
        return $this->belongsTo('AdminModel', 'admin_id', 'id');
    }

    /**
     * 文章添加
     * @param $data
     * @return array|int|string
     */
    public function add($data)
    {
        $validate = new ArticleValidate();
        if (!$validate->scene('add')->check($data)){
            return $validate->getError();
        }

        $result = $this->save($data);
        if (!$result) {
            return "服务器错误,请稍后再试！";
        }
        return 1;
    }
}
