<?php


namespace app\common\controller\admin;


use think\Controller;

class AdminController extends Controller
{

    /**
     * 通用删除方法
     * @param $model
     */
    public function delete($model)
    {
        $id = request()->post('id');
        $result = $model->useSoftDelete('delete_time', time())->delete($id);

        if (!$result)
            $this->error("删除失败");
        else
            $this->success("删除成功！");
    }


    public function change($model, $field, $value='value')
    {
        $value = request()->post($value, 1);
        $result = $model
            ->where('id', request()->post('id'))
            ->update([$field => $value]);
        if (!$result)
            $this->error("修改失败");
        else
            $this->success("修改成功！");
    }

}