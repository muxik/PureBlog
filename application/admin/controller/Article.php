<?php


namespace app\admin\controller;


use app\common\model\ArticleModel;
use think\Controller;
use think\Image;
use think\Request;

class Article extends Controller
{

    protected $base = ['source', 'pic', 'content', 'title', 'category_id', 'tag', 'description'];
    protected $model;

    protected function initialize()
    {
        $this->model = new ArticleModel();
        parent::initialize();
    }


    public function index()
    {
        $article = $this->model;

        return view()->assign([
            "count" => $article->count(),
            "article" => $article
                ->field('id,title,pic,state,read,description,tag,top,create_time,category_id,admin_id')
                ->with(['category', 'admin'])
                ->paginate(10),
        ]);
    }

    /**
     * 文章创作
     * @return \think\response\View
     */
    public function create()
    {
        return view()
            ->assign([
                "category" => model('CategoryModel')
                    ->where('type', 'neq', 1)
                    ->select()
            ]);
    }


    /**
     * 文章修改
     * @param $id
     * @return \think\response\View
     */
    public function edit($id)
    {
        return view()
            ->assign([
                "category" => model('CategoryModel')
                    ->where('type', 'neq', 1)
                    ->select(),
                'article' => $this
                    ->model
                    ->field('id,source,lock,title,pic,state,read,description,tag,top,create_time,category_id')
                    ->where('id', $id)
                    ->with(['category'])->find()
            ]);
    }


    /**
     * 更改文章状态
     * @param Request $request
     */
    public function changeState(Request $request)
    {
        $state = $request->post('value', 1);
        $result = $this
            ->model
            ->where('id', $request->post('id'))
            ->update(['state' => $state]);
        if (!$result)
            $this->error("修改失败");
        else
            $this->success("修改成功！");
    }

    /**
     * 更改文章置顶
     * @param Request $request
     */
    public function changeTop(Request $request)
    {
        $state = $request->post('value', 1);
        $result = $this
            ->model
            ->where('id', $request->post('id'))
            ->update(['top' => $state]);
        if (!$result)
            $this->error("修改失败");
        else
            $this->success("修改成功！");
    }

    /**
     * 文章添加
     * @param Request $request
     */
    public function add(Request $request)
    {
        // 如果文章设置了加密就接收lock & key
        if (!empty($request->post('lock')))
            array_push($this->base, 'lock','key');


        $data = $request->only($this->base);

        session(['prefix' => 'admin_',]);
        $id = session('user.id');
        $data['admin_id'] = $id != null ? $id : cookie('admin_user')['id'];

        $result = $this->model->add($data);

        if (1 !== $result) {
            $this->error($result);
        }
        $this->success('添加成功！');
    }


    public function update(Request $request)
    {

        // 如果文章设置了加密就接收lock & key
        if (!empty($request->post('lock')))
            array_push($this->base, 'lock', 'key');


        $data = $request->only($this->base);
        $data['id'] = $request->post('id');
        session(['prefix' => 'admin_',]);
        $id = session('user.id');
        $data['admin_id'] = $id != null ? $id : cookie('admin_user')['id'];

        $result = $this->model->edit($data);

        if (1 !== $result) {
            $this->error($result);
        }

        $this->success('修改成功');
    }

    /**
     * 文章封面上传
     * @param Request $request
     * @return \think\response\Json
     */
    public function upload(Request $request)
    {
        $file = $request->file('file');
        $info = $file
            ->validate(['size' => 15678888888, 'ext' => 'jpg,png,gif'])
            ->move('./static/upload');

        $save_path = '/static/upload/' . $info->getSaveName();

        $image = Image::open('.' . $save_path);
        $image->thumb(1280, 1280)->save('.' . $save_path);

        if ($info) {
            return json(['code' => 1, 'files' => ['file' => $save_path]]);
        }
    }

    /**
     * 删除分类
     * @param Request $request
     * @throws \Exception
     */
    public function del(Request $request)
    {
        $id = $request->post('id');
        $result = $this->model->useSoftDelete('delete_time', time())->delete($id);

        if (!$result)
            $this->error("删除失败");
        else
            $this->success("删除成功！");
    }

}