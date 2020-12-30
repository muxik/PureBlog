<?php


namespace app\admin\controller;


use app\admin\model\ArticleModel;
use think\Controller;
use think\Image;
use think\Request;

class Article extends Controller
{

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
            ->assign(["category" => model('CategoryModel')->where('type', 'neq', 1)->select(),]);
    }


    public function add(Request $request)
    {

        $base = ['pic', 'content', 'title', 'category_id', 'tag', 'description'];

        // 如果文章设置了加密就接收lock & key
        if (!empty($request->post('lock')))
            array_push($base, 'lock');
        array_push($base, 'key');

        $data = $request->only($base);

        session(['prefix' => 'admin_',]);
        $id = session('user.id');
        $data['admin_id'] = $id != null ? $id : cookie('admin_user')['id'];

        $result = $this->model->add($data);

        if (1 !== $result) {
            $this->error($result);
        }
        $this->success('添加成功！');
    }

    /**
     * 文章封面上传
     * @param Request $request
     * @return \think\response\Json
     */
    public function upload(Request $request)
    {
        $file = $request->file('file');
        $info = $file->validate(['size' => 15678888888, 'ext' => 'jpg,png,gif'])->move('./static/upload');

        $save_path = '/static/upload/' . $info->getSaveName();

        $image = Image::open('.' . $save_path);
        $image->thumb(1280, 1280)->save('.' . $save_path);

        if ($info) {
            return json(['code' => 1, 'files' => ['file' => $save_path]]);
        }
    }

}