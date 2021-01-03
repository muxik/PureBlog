<?php


namespace app\common\model;



use app\common\validate\LinkValidate;
use think\Model;
use think\model\concern\SoftDelete;

class LinkModel extends Model
{

    use SoftDelete;

    protected $name = 'link';

    /**
     * 友情链接添加
     * @param $info
     * @return array|int|string
     */
    public function add($info)
    {
        $validate = new LinkValidate();
        if (!$validate->scene('add')->check($info)) {
            return $validate->getError();
        }

        $result = $this->save($info);
        if (!$result) {
            return '服务器错误,请稍后再试!';
        }

        return 1;

    }

    /**
     * 友情链接编辑
     * @param $info
     * @return array|int|string
     */
    public function edit($info)
    {
        $validate = new LinkValidate();
        if (!$validate->scene('add')->check($info)) {
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