<?php


namespace app\admin\model;


use app\admin\validate\LinkValidate;
use think\Model;
use think\model\concern\SoftDelete;

class LinkModel extends Model
{

    use SoftDelete;

    protected $name = 'link';

    public function add($info)
    {
        $validate = new LinkValidate();
        if (!$validate->scene('add')->check($info)){
            return $validate->getError();
        }

        $result = $this->save($info);
        if (!$result){
            return  '服务器错误,请稍后再试!';
        }

        return 1;

    }
}