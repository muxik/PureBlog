<?php

namespace app\common\model;

use think\Model;
use think\model\concern\SoftDelete;

class Admin extends Model
{
    use SoftDelete;

    protected $name = "admin";

    function login()
    {

    }
}
