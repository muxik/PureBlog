<?php


namespace app\common\model;



use app\common\validate\CategoryValidate;
use think\Model;
use think\model\concern\SoftDelete;

class CategoryModel extends Model
{

    use SoftDelete;

    protected $name = 'category';


    public function add($data)
    {
        $validate = new CategoryValidate();
        if (!$validate->scene('add')->check($data)){
            return $validate->getError();
        }

        $result = $this->save($data);
        if (!$result){
            return "服务器错误，请稍后再试！";
        }
        return 1;
    }


    public function edit($data)
    {
        $validate = new CategoryValidate();
        if (!$validate->scene('edit')->check($data)){
            return $validate->getError();
        }

        $result = $this
            ->where('id', $data['id'])
            ->update($data);

        if (!$result){
            return '服务器错误！请稍后再试';
        }

        return 1;
    }

    /**
     * 关联文章
     * @return \think\model\relation\HasMany
     */
    public function article()
    {
        return $this->hasMany('ArticleModel', 'category_id', 'id');
    }
}