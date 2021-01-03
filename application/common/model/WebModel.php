<?php


namespace app\common\model;


use app\common\validate\WebValidate;
use think\Model;
use think\model\concern\SoftDelete;

class WebModel extends Model
{

    use SoftDelete;

    protected $name = 'web';


    /**
     * web info编辑
     * @param $info
     * @return array|int|string
     */
    public function edit($info)
    {
        $validate = new WebValidate();
        if (!$validate->scene('edit')->check($info)) {
            return $validate->getError();
        }

        $id = $info['id'];
        unset($info['id']);
        $result = $this->where('id', $id)->update($info);
        if (!$result) {
            return '服务器错误,请稍后再试!';
        }
        return 1;

    }
}