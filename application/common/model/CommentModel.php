<?php


namespace app\common\model;


use app\common\validate\CommentValidate;
use think\Model;
use think\model\concern\SoftDelete;

class CommentModel extends Model
{
    protected $name = "comment";

    use SoftDelete;

    public function add($info)
    {
        $validate = new CommentValidate();
        if (!$validate->scene('add')->check($info)){
            return $validate->getError();
        }

        $result = $this->save($info);
        if (!$result){
            return "服务器错误，请稍后再试！";
        }
        return 1;
    }


    public function scomment()
    {
        $this->hasMany('CommentModel', 'pid','id');
    }

}