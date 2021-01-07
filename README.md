# PureBlog v2.0.9

基于thinkphp5.1风格简约功能丰富的php博客系统

## 我的运行环境

> PHP >= 5.6

- 环境
    - ArchLinux
    - PHP/7.4
    - Nginx/1.18
- Composer/1.10.16
- PHP扩展
    - PDO PHP Extension
    - MBstring PHP Extension
    - GD PHP Extension

---



## MySQL 配置

####  连接MySQL

修改 `config/database.php`

#### 导入数据库文件

`blog.sql`



---



## URL重写

#### 1、[Apache]

1. `httpd.conf`配置文件中加载了`mod_rewrite.so`模块
2. `AllowOverride None` 将`None`改为 `All`
3. 把下面的内容保存为`.htaccess`文件放到应用入口文件的同级目录下

```
<IfModule mod_rewrite.c>
  Options +FollowSymlinks -Multiviews
  RewriteEngine On

  RewriteCond %{REQUEST_FILENAME} !-d
  RewriteCond %{REQUEST_FILENAME} !-f
  RewriteRule ^(.*)$ index.php/$1 [QSA,PT,L]
</IfModule>
```

#### 2、[Nginx]

```
try_files $uri $uri/ /index.php?s=$uri&$args;  
```

Nginx低版本

```
location / { // …..省略部分代码
   if (!-e $request_filename) {
   		rewrite  ^(.*)$  /index.php?s=/$1  last;
    }
}
```

> 注意: 设置url重写后静态资源路径会失效，可以通过一下方式设置


```
location /static/ {
	root /srv/nginx/blog/public/;
	# 请根据实际情况修改
}

```

## 更新日至
  - 2.0.9 重写所有功能新增文章评论多级分类

## LICENSE

Apache