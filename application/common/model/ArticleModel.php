<?php


namespace app\common\model;


use app\common\validate\ArticleValidate;
use think\Model;
use think\model\concern\SoftDelete;
use think\model\relation\BelongsTo;

class ArticleModel extends Model
{
    use SoftDelete;

    protected $name = "article";


    /**
     * 关联栏目
     * @return BelongsTo
     */
    public function category()
    {
        return $this->belongsTo('CategoryModel', 'category_id', 'id');
    }

    /**
     * 关联作者
     * @return BelongsTo
     */
    public function admin()
    {
        return $this->belongsTo('app\\admin\\model\\AdminModel', 'admin_id', 'id');
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

    public function edit($data)
    {
        $validate = new ArticleValidate();
        if (!$validate->scene('add')->check($data)){
            return $validate->getError();
        }
        $id = $data['id'];
        unset($data['id']);

        $result = $this->where('id', $id)->update($data);
        if (!$result) {
            return "服务器错误,请稍后再试！";
        }
        return 1;
    }

    /**
     * @param $keyword
     * @param $value
     * @return ArticleModel
     */
    public function search($keyword, $value)
    {
        $where[] = [
            $value, 'like', '%' . $keyword . '%'
        ];
        return $this->where($where);
    }


    /**
     * 关联评论
     * @return ArticleModel|\think\model\relation\HasMany
     */
    public function comments()
    {
        return $this->hasMany('CommentModel', 'article_id', 'id');
    }

}
