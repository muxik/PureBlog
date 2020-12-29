<?php


namespace app\admin\model;


use app\admin\validate\CategoryValidate;
use think\Model;

class CategoryModel extends Model
{
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
}