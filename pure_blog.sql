/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80021
 Source Host           : localhost:3306
 Source Schema         : pure_blog

 Target Server Type    : MySQL
 Target Server Version : 80021
 File Encoding         : 65001

 Date: 11/01/2021 11:47:55
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_admin
-- ----------------------------
DROP TABLE IF EXISTS `t_admin`;
CREATE TABLE `t_admin` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `password` varchar(255) NOT NULL,
  `state` tinyint(1) NOT NULL DEFAULT '1',
  `email` varchar(255) NOT NULL,
  `nickname` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `create_time` int NOT NULL,
  `update_time` int NOT NULL,
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_admin
-- ----------------------------
BEGIN;
INSERT INTO `t_admin` VALUES (1, 'root', '290421b48e4430573fa15d5a9d4f42f4', 2, 'lqjxm666@163.com', 'Muxi_k', 1609330779, 1609330779, NULL);
INSERT INTO `t_admin` VALUES (2, 'test', '4fd0ae9f6c7638ac36e4d55585a78881', 0, 'Test@gmail.com', 'Test', 1609405272, 1609405272, NULL);
COMMIT;

-- ----------------------------
-- Table structure for t_article
-- ----------------------------
DROP TABLE IF EXISTS `t_article`;
CREATE TABLE `t_article` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `source` text,
  `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `admin_id` int NOT NULL,
  `description` text NOT NULL,
  `pic` varchar(255) NOT NULL,
  `category_id` int NOT NULL,
  `read` int DEFAULT '0',
  `key` varchar(255) DEFAULT NULL,
  `lock` tinyint(1) DEFAULT '0',
  `state` tinyint(1) NOT NULL DEFAULT '1',
  `top` tinyint(1) DEFAULT '0',
  `tag` varchar(255) NOT NULL,
  `create_time` int NOT NULL,
  `update_time` int NOT NULL,
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_article
-- ----------------------------
BEGIN;
INSERT INTO `t_article` VALUES (5, '霉运的一天', '以前记录下来的，迟迟没有发出来。\n\ncreation time, 2019.8.30 00:12\n\n昨晚下过雨，早晨拿出纸巾擦好自行车座位然后扫码，该车已故障，突然想起来这辆车自己昨天扫的时候就已经故障了还停在这里\n\n去学校的半路上突然下起了雨，慢慢变大了\n\n学校门前的最后一个红绿灯错以为时间到了下一刻就是绿灯可以通行，结果不是，侧边的车蜂拥一般袭来，还好躲得快\n\n到理科楼下停下来的时候发现自行车脚撑是坏的，弄了好几次才勉勉强强立起来\n\n中午吃过饭打算去 WAIC，可是依然在下着雨，一手撑伞一手骑车真的没有尝试过，过了一个减速带突然控制不住方向撞在了侧边，鞋面磨破了一个洞\n\n停车的时候弄倒了旁边的车\n\n进地铁的时候拿着的伞突然不知道怎么了夹在了门口的挡板上，由于惯性折断了一个角上的塑料连接\n\n去世博馆的路上由于导航以及错误的指导走错了入口无法进入只好折路而返走另一边的门\n\n无人驾驶因为下雨并没有出行\n\n回去的路上刚到路口绿灯变成了红灯\n\n出地铁站扫一辆车打算回到学校，和早晨类似的经历，刚擦好座位扫码开锁推了出来，然后发现没有脚踏只好重新扫一辆\n\n晚饭点了一份没有吃过的菜超级辣\n\n去将餐具送回去的座位间的小道上因为空调的风吹了过来，WAIC 会议拿到的口袋被固定在了座位一边然后没有反应过来撞了上去，手里的餐具差点滑掉，但是口袋还是折掉了\n\n晚上回去的路上好不容易找到的小蓝车，擦好座位没想到突然倒了过来撞在了身上，手中的手机也顺势掉在了地上，无奈这辆车依然无法开锁，关键的是平台显示锁已经打开了，反馈以后无法再次扫车，一直显示正在处理中，于是找了一个没有月卡的摩拜骑了回来\n\n回到住的地方却发现楼下二十多天从来没有锁过的门突然锁掉了，只好联系老板开门\n\nQAQ 希望明天会幸运一点点', '<p>以前记录下来的，迟迟没有发出来。</p>\n<p>creation time, 2019.8.30 00:12</p>\n<p>昨晚下过雨，早晨拿出纸巾擦好自行车座位然后扫码，该车已故障，突然想起来这辆车自己昨天扫的时候就已经故障了还停在这里</p>\n<p>去学校的半路上突然下起了雨，慢慢变大了</p>\n<p>学校门前的最后一个红绿灯错以为时间到了下一刻就是绿灯可以通行，结果不是，侧边的车蜂拥一般袭来，还好躲得快</p>\n<p>到理科楼下停下来的时候发现自行车脚撑是坏的，弄了好几次才勉勉强强立起来</p>\n<p>中午吃过饭打算去 WAIC，可是依然在下着雨，一手撑伞一手骑车真的没有尝试过，过了一个减速带突然控制不住方向撞在了侧边，鞋面磨破了一个洞</p>\n<p>停车的时候弄倒了旁边的车</p>\n<p>进地铁的时候拿着的伞突然不知道怎么了夹在了门口的挡板上，由于惯性折断了一个角上的塑料连接</p>\n<p>去世博馆的路上由于导航以及错误的指导走错了入口无法进入只好折路而返走另一边的门</p>\n<p>无人驾驶因为下雨并没有出行</p>\n<p>回去的路上刚到路口绿灯变成了红灯</p>\n<p>出地铁站扫一辆车打算回到学校，和早晨类似的经历，刚擦好座位扫码开锁推了出来，然后发现没有脚踏只好重新扫一辆</p>\n<p>晚饭点了一份没有吃过的菜超级辣</p>\n<p>去将餐具送回去的座位间的小道上因为空调的风吹了过来，WAIC 会议拿到的口袋被固定在了座位一边然后没有反应过来撞了上去，手里的餐具差点滑掉，但是口袋还是折掉了</p>\n<p>晚上回去的路上好不容易找到的小蓝车，擦好座位没想到突然倒了过来撞在了身上，手中的手机也顺势掉在了地上，无奈这辆车依然无法开锁，关键的是平台显示锁已经打开了，反馈以后无法再次扫车，一直显示正在处理中，于是找了一个没有月卡的摩拜骑了回来</p>\n<p>回到住的地方却发现楼下二十多天从来没有锁过的门突然锁掉了，只好联系老板开门</p>\n<p>QAQ 希望明天会幸运一点点</p>\n', 1, '就像上天专想与你作对似的，一整天都是糟糕的~\n<br>\n希望明天会幸运一点点~\n\n', '/static/upload/20210103/008ef32b33896fee20c5c66e04cae98d.png', 3, 0, '', 1, 1, 1, '心情,记录', 1609330779, 1610000504, 1610000504);
INSERT INTO `t_article` VALUES (6, '你好', '\n[TOCM]\n\n## 1、同步和异步的区别\n\n- **同步**: 所有的操作都做完，才返回给用户。这样用户在线等待的时间太长，给用户一种卡死了的感觉\n\n- **异步**: 不等所有操作做完 ，就返回给用户。操作已经启动程序会将它加到消息队列当中，慢慢执行，不会给用户卡死的感觉\n\n\n\n## 2、使用Asynico 进行异步编程\n\n- Asynico  - 异步 I/O\n\n  - asyncio 是用来编写 **并发** 代码的库，使用 **async/await** 语法。\n\n  - asyncio 被用作多个提供高性能 Python 异步框架的基础，包括网络和网站服务，数据库连接库，分布式任务队列等等。\n\n  - asyncio 往往是构建 IO 密集型和高层级 **结构化** 网络代码的最佳选择。\n\n> python 3.5 后才可以使用async/await 语法\n\n\n\n### 异步请求图片示例：\n\n\n\n安装异步http 请求库\n\n```sh\n$ pip install aiohttp\n```\n\n\n\n- 导入 asyncio 以及异步请求aiohttp\n\n ```python\n  # !/usr/bin/env python\n  \n  import asyncio\n  import aiohttp\n  ```\n\n- fetch \n\n ```python\n  async def fetch(session, url):\n      async with session.get(url, headers={\n          \"referer\": \"https://www.mzitu.com/\",\n          \"user-agent\": \"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36\"\n          }) as response:\n          return await response.read()\n  \n  ```\n\n-  download_image()\n\n  ```python\n  async def download_image(url):\n      async with aiohttp.ClientSession() as session:\n          result = await fetch(session, url)\n          with open(url.split(\"/\")[-1], \'wb\') as f:\n              f.write(result)\n          f.close()\n   ```\n\n-  __main()__\n\n  ```python\n  if __name__ == \'__main__\':\n      \n      # 添加多个任务\n      tasks = [\n              asyncio.ensure_future(download_image(\"https://csdnimg.cn/cdn/content-toolbar/csdn-logo.png?v=20200416.1\")),\n              asyncio.ensure_future(download_image(\"https://csdnimg.cn/public/common/toolbar/images/csdnqr@2x.png\"))\n              ]\n      # 打开时间循环\n      loop = asyncio.get_event_loop()\n      \n      loop.run_until_complete(asyncio.gather(*tasks))\n  \n  ```\n\n\n\n\n\n## 参考\n\n[https://docs.python.org/zh-cn/3/library/asyncio.html](https://docs.python.org/zh-cn/3/library/asyncio.html)\n\n[https://docs.aiohttp.org/en/stable/](https://docs.aiohttp.org/en/stable/)\n\n', '<p>[TOCM]</p>\n<h2 id=\"h2-1-\"><a name=\"1、同步和异步的区别\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>1、同步和异步的区别</h2><ul>\n<li><p><strong>同步</strong>: 所有的操作都做完，才返回给用户。这样用户在线等待的时间太长，给用户一种卡死了的感觉</p>\n</li><li><p><strong>异步</strong>: 不等所有操作做完 ，就返回给用户。操作已经启动程序会将它加到消息队列当中，慢慢执行，不会给用户卡死的感觉</p>\n</li></ul>\n<h2 id=\"h2-2-asynico-\"><a name=\"2、使用Asynico 进行异步编程\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>2、使用Asynico 进行异步编程</h2><ul>\n<li><p>Asynico  - 异步 I/O</p>\n<ul>\n<li><p>asyncio 是用来编写 <strong>并发</strong> 代码的库，使用 <strong>async/await</strong> 语法。</p>\n</li><li><p>asyncio 被用作多个提供高性能 Python 异步框架的基础，包括网络和网站服务，数据库连接库，分布式任务队列等等。</p>\n</li><li><p>asyncio 往往是构建 IO 密集型和高层级 <strong>结构化</strong> 网络代码的最佳选择。</p>\n</li></ul>\n</li></ul>\n<blockquote>\n<p>python 3.5 后才可以使用async/await 语法</p>\n</blockquote>\n<h3 id=\"h3--\"><a name=\"异步请求图片示例：\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>异步请求图片示例：</h3><p>安装异步http 请求库</p>\n<pre class=\"prettyprint linenums prettyprinted\" style=\"\"><ol class=\"linenums\"><li class=\"L0\"><code class=\"lang-sh\"><span class=\"pln\">$ pip install aiohttp</span></code></li></ol></pre>\n<ul>\n<li><p>导入 asyncio 以及异步请求aiohttp</p>\n<pre class=\"prettyprint linenums prettyprinted\" style=\"\"><ol class=\"linenums\"><li class=\"L0\"><code class=\"lang-python\"><span class=\"com\"># !/usr/bin/env python</span></code></li><li class=\"L1\"><code class=\"lang-python\"></code></li><li class=\"L2\"><code class=\"lang-python\"><span class=\"kwd\">import</span><span class=\"pln\"> asyncio</span></code></li><li class=\"L3\"><code class=\"lang-python\"><span class=\"kwd\">import</span><span class=\"pln\"> aiohttp</span></code></li></ol></pre>\n</li><li><p>fetch </p>\n<pre class=\"prettyprint linenums prettyprinted\" style=\"\"><ol class=\"linenums\"><li class=\"L0\"><code class=\"lang-python\"><span class=\"pln\">async </span><span class=\"kwd\">def</span><span class=\"pln\"> fetch</span><span class=\"pun\">(</span><span class=\"pln\">session</span><span class=\"pun\">,</span><span class=\"pln\"> url</span><span class=\"pun\">):</span></code></li><li class=\"L1\"><code class=\"lang-python\"><span class=\"pln\">    async </span><span class=\"kwd\">with</span><span class=\"pln\"> session</span><span class=\"pun\">.</span><span class=\"pln\">get</span><span class=\"pun\">(</span><span class=\"pln\">url</span><span class=\"pun\">,</span><span class=\"pln\"> headers</span><span class=\"pun\">={</span></code></li><li class=\"L2\"><code class=\"lang-python\"><span class=\"pln\">        </span><span class=\"str\">\"referer\"</span><span class=\"pun\">:</span><span class=\"pln\"> </span><span class=\"str\">\"https://www.mzitu.com/\"</span><span class=\"pun\">,</span></code></li><li class=\"L3\"><code class=\"lang-python\"><span class=\"pln\">        </span><span class=\"str\">\"user-agent\"</span><span class=\"pun\">:</span><span class=\"pln\"> </span><span class=\"str\">\"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36\"</span></code></li><li class=\"L4\"><code class=\"lang-python\"><span class=\"pln\">        </span><span class=\"pun\">})</span><span class=\"pln\"> </span><span class=\"kwd\">as</span><span class=\"pln\"> response</span><span class=\"pun\">:</span></code></li><li class=\"L5\"><code class=\"lang-python\"><span class=\"pln\">        </span><span class=\"kwd\">return</span><span class=\"pln\"> await response</span><span class=\"pun\">.</span><span class=\"pln\">read</span><span class=\"pun\">()</span></code></li></ol></pre>\n</li><li><p>download_image()</p>\n<pre class=\"prettyprint linenums prettyprinted\" style=\"\"><ol class=\"linenums\"><li class=\"L0\"><code class=\"lang-python\"><span class=\"pln\">async </span><span class=\"kwd\">def</span><span class=\"pln\"> download_image</span><span class=\"pun\">(</span><span class=\"pln\">url</span><span class=\"pun\">):</span></code></li><li class=\"L1\"><code class=\"lang-python\"><span class=\"pln\">   async </span><span class=\"kwd\">with</span><span class=\"pln\"> aiohttp</span><span class=\"pun\">.</span><span class=\"typ\">ClientSession</span><span class=\"pun\">()</span><span class=\"pln\"> </span><span class=\"kwd\">as</span><span class=\"pln\"> session</span><span class=\"pun\">:</span></code></li><li class=\"L2\"><code class=\"lang-python\"><span class=\"pln\">       result </span><span class=\"pun\">=</span><span class=\"pln\"> await fetch</span><span class=\"pun\">(</span><span class=\"pln\">session</span><span class=\"pun\">,</span><span class=\"pln\"> url</span><span class=\"pun\">)</span></code></li><li class=\"L3\"><code class=\"lang-python\"><span class=\"pln\">       </span><span class=\"kwd\">with</span><span class=\"pln\"> open</span><span class=\"pun\">(</span><span class=\"pln\">url</span><span class=\"pun\">.</span><span class=\"pln\">split</span><span class=\"pun\">(</span><span class=\"str\">\"/\"</span><span class=\"pun\">)[-</span><span class=\"lit\">1</span><span class=\"pun\">],</span><span class=\"pln\"> </span><span class=\"str\">\'wb\'</span><span class=\"pun\">)</span><span class=\"pln\"> </span><span class=\"kwd\">as</span><span class=\"pln\"> f</span><span class=\"pun\">:</span></code></li><li class=\"L4\"><code class=\"lang-python\"><span class=\"pln\">           f</span><span class=\"pun\">.</span><span class=\"pln\">write</span><span class=\"pun\">(</span><span class=\"pln\">result</span><span class=\"pun\">)</span></code></li><li class=\"L5\"><code class=\"lang-python\"><span class=\"pln\">       f</span><span class=\"pun\">.</span><span class=\"pln\">close</span><span class=\"pun\">()</span></code></li></ol></pre>\n</li><li><p><strong>main()</strong></p>\n<pre class=\"prettyprint linenums prettyprinted\" style=\"\"><ol class=\"linenums\"><li class=\"L0\"><code class=\"lang-python\"><span class=\"kwd\">if</span><span class=\"pln\"> __name__ </span><span class=\"pun\">==</span><span class=\"pln\"> </span><span class=\"str\">\'__main__\'</span><span class=\"pun\">:</span></code></li><li class=\"L1\"><code class=\"lang-python\"></code></li><li class=\"L2\"><code class=\"lang-python\"><span class=\"pln\">   </span><span class=\"com\"># 添加多个任务</span></code></li><li class=\"L3\"><code class=\"lang-python\"><span class=\"pln\">   tasks </span><span class=\"pun\">=</span><span class=\"pln\"> </span><span class=\"pun\">[</span></code></li><li class=\"L4\"><code class=\"lang-python\"><span class=\"pln\">           asyncio</span><span class=\"pun\">.</span><span class=\"pln\">ensure_future</span><span class=\"pun\">(</span><span class=\"pln\">download_image</span><span class=\"pun\">(</span><span class=\"str\">\"https://csdnimg.cn/cdn/content-toolbar/csdn-logo.png?v=20200416.1\"</span><span class=\"pun\">)),</span></code></li><li class=\"L5\"><code class=\"lang-python\"><span class=\"pln\">           asyncio</span><span class=\"pun\">.</span><span class=\"pln\">ensure_future</span><span class=\"pun\">(</span><span class=\"pln\">download_image</span><span class=\"pun\">(</span><span class=\"str\">\"https://csdnimg.cn/public/common/toolbar/images/csdnqr@2x.png\"</span><span class=\"pun\">))</span></code></li><li class=\"L6\"><code class=\"lang-python\"><span class=\"pln\">           </span><span class=\"pun\">]</span></code></li><li class=\"L7\"><code class=\"lang-python\"><span class=\"pln\">   </span><span class=\"com\"># 打开时间循环</span></code></li><li class=\"L8\"><code class=\"lang-python\"><span class=\"pln\">   loop </span><span class=\"pun\">=</span><span class=\"pln\"> asyncio</span><span class=\"pun\">.</span><span class=\"pln\">get_event_loop</span><span class=\"pun\">()</span></code></li><li class=\"L9\"><code class=\"lang-python\"></code></li><li class=\"L0\"><code class=\"lang-python\"><span class=\"pln\">   loop</span><span class=\"pun\">.</span><span class=\"pln\">run_until_complete</span><span class=\"pun\">(</span><span class=\"pln\">asyncio</span><span class=\"pun\">.</span><span class=\"pln\">gather</span><span class=\"pun\">(*</span><span class=\"pln\">tasks</span><span class=\"pun\">))</span></code></li></ol></pre>\n</li></ul>\n<h2 id=\"h2-u53C2u8003\"><a name=\"参考\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>参考</h2><p><a href=\"https://docs.python.org/zh-cn/3/library/asyncio.html\">https://docs.python.org/zh-cn/3/library/asyncio.html</a></p>\n<p><a href=\"https://docs.aiohttp.org/en/stable/\">https://docs.aiohttp.org/en/stable/</a></p>\n', 1, '伤口上的考古发掘', '/static/upload/20210104/34c5353ff729b8ff092771273543b858.jpg', 2, 0, NULL, 0, 1, 0, '的是看风景', 1609332380, 1610000533, 1610000533);
INSERT INTO `t_article` VALUES (7, 'Vue cli3 跨域请求', '\n# Vue cli3 跨域请求\n ---\n\n## 问题\n\n最近学了Vue 想写点demo练习一下 就用tp写了个接口\n\n - 此时\n 	+ server：api.org:80/users\n    + client : 127.0.0.1:8080\n\n当我写完请求后发现出问题了\n```javascript\n// 请求代码\ngetUsersTable() {\n      this.$http\n        .get(\"api.org/users\")\n        .then(res => {\n          console.log(res);\n          this.users = res.body;\n        })\n        .catch(err => {\n          console.error(err);\n        });\n    }\n```\n__报错如下：__\n   `Access to XMLHttpRequest at \'http://api.org/users\' from origin \'http://127.0.0.1:8080\' has been blocked by CORS policy: No \'Access-Control-Allow-Origin\' header is present on the requested resource.`\n   \n   发现是跨域问题\n\n## 解决方法　\n1. 配置文件\n	 创建配置文件`vue.config.js`,与`package.json` 同级\n2. 安装代理\n    现在vue cli3默认是没有devServer的，需要手动安装\n ```bash\n     $ npm install --save-dev http-proxy-middleware\n```\n\n3. 写入配置\n\n```javascript\nmodule.exports = {\n    devServer: {\n        proxy: {\n            //配置跨域\n            \'/api\': {\n                target: \"http://api.org\",\n                ws: true,\n                changOrigin: true,\n                pathRewrite: {\n                    \'^/api\': \'/\'\n                }\n            }\n        }\n    }\n}\n```\n这里要把请求路径该一下：\n`.get(\"api.org/users\")`\n改成:\n`.get(\"api/users\")`\n\n### 注意问题\n如果你是用tp写的接口\n把	`app.config里的　app_debug = true 改为　false`\n**linux下**\n*runtime 权限改为777*\n```bash\n$ sudo chmod 777 runtime \n```\n__否则请求的时候会有服务器500错误__\n\n--- \n## 参考文献\n__[https://cli.vuejs.org/zh/config/#devserver-proxy](https://cli.vuejs.org/zh/config/#devserver-proxy)__\n__[https://github.com/chimurai/http-proxy-middleware#proxycontext-config](https://github.com/chimurai/http-proxy-middleware#proxycontext-config)__\n', '<h1 id=\"h1-vue-cli3-\"><a name=\"Vue cli3 跨域请求\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>Vue cli3 跨域请求</h1><hr>\n<h2 id=\"h2-u95EEu9898\"><a name=\"问题\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>问题</h2><p>最近学了Vue 想写点demo练习一下 就用tp写了个接口</p>\n<ul>\n<li>此时<ul>\n<li>server：api.org:80/users<ul>\n<li>client : 127.0.0.1:8080</li></ul>\n</li></ul>\n</li></ul>\n<p>当我写完请求后发现出问题了</p>\n<pre class=\"prettyprint linenums prettyprinted\" style=\"\"><ol class=\"linenums\"><li class=\"L0\"><code class=\"lang-javascript\"><span class=\"com\">// 请求代码</span></code></li><li class=\"L1\"><code class=\"lang-javascript\"><span class=\"pln\">getUsersTable</span><span class=\"pun\">()</span><span class=\"pln\"> </span><span class=\"pun\">{</span></code></li><li class=\"L2\"><code class=\"lang-javascript\"><span class=\"pln\">      </span><span class=\"kwd\">this</span><span class=\"pun\">.</span><span class=\"pln\">$http</span></code></li><li class=\"L3\"><code class=\"lang-javascript\"><span class=\"pln\">        </span><span class=\"pun\">.</span><span class=\"kwd\">get</span><span class=\"pun\">(</span><span class=\"str\">\"api.org/users\"</span><span class=\"pun\">)</span></code></li><li class=\"L4\"><code class=\"lang-javascript\"><span class=\"pln\">        </span><span class=\"pun\">.</span><span class=\"pln\">then</span><span class=\"pun\">(</span><span class=\"pln\">res </span><span class=\"pun\">=&gt;</span><span class=\"pln\"> </span><span class=\"pun\">{</span></code></li><li class=\"L5\"><code class=\"lang-javascript\"><span class=\"pln\">          console</span><span class=\"pun\">.</span><span class=\"pln\">log</span><span class=\"pun\">(</span><span class=\"pln\">res</span><span class=\"pun\">);</span></code></li><li class=\"L6\"><code class=\"lang-javascript\"><span class=\"pln\">          </span><span class=\"kwd\">this</span><span class=\"pun\">.</span><span class=\"pln\">users </span><span class=\"pun\">=</span><span class=\"pln\"> res</span><span class=\"pun\">.</span><span class=\"pln\">body</span><span class=\"pun\">;</span></code></li><li class=\"L7\"><code class=\"lang-javascript\"><span class=\"pln\">        </span><span class=\"pun\">})</span></code></li><li class=\"L8\"><code class=\"lang-javascript\"><span class=\"pln\">        </span><span class=\"pun\">.</span><span class=\"kwd\">catch</span><span class=\"pun\">(</span><span class=\"pln\">err </span><span class=\"pun\">=&gt;</span><span class=\"pln\"> </span><span class=\"pun\">{</span></code></li><li class=\"L9\"><code class=\"lang-javascript\"><span class=\"pln\">          console</span><span class=\"pun\">.</span><span class=\"pln\">error</span><span class=\"pun\">(</span><span class=\"pln\">err</span><span class=\"pun\">);</span></code></li><li class=\"L0\"><code class=\"lang-javascript\"><span class=\"pln\">        </span><span class=\"pun\">});</span></code></li><li class=\"L1\"><code class=\"lang-javascript\"><span class=\"pln\">    </span><span class=\"pun\">}</span></code></li></ol></pre>\n<p><strong>报错如下：</strong><br>   <code>Access to XMLHttpRequest at \'http://api.org/users\' from origin \'http://127.0.0.1:8080\' has been blocked by CORS policy: No \'Access-Control-Allow-Origin\' header is present on the requested resource.</code></p>\n<p>   发现是跨域问题</p>\n<h2 id=\"h2-u89E3u51B3u65B9u6CD5\"><a name=\"解决方法\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>解决方法　</h2><ol>\n<li>配置文件<br>  创建配置文件<code>vue.config.js</code>,与<code>package.json</code> 同级</li><li><p>安装代理<br> 现在vue cli3默认是没有devServer的，需要手动安装</p>\n<pre class=\"prettyprint linenums prettyprinted\" style=\"\"><ol class=\"linenums\"><li class=\"L0\"><code class=\"lang-bash\"><span class=\"pln\">  $ npm install </span><span class=\"pun\">--</span><span class=\"pln\">save</span><span class=\"pun\">-</span><span class=\"pln\">dev http</span><span class=\"pun\">-</span><span class=\"pln\">proxy</span><span class=\"pun\">-</span><span class=\"pln\">middleware</span></code></li></ol></pre>\n</li><li><p>写入配置</p>\n</li></ol>\n<pre class=\"prettyprint linenums prettyprinted\" style=\"\"><ol class=\"linenums\"><li class=\"L0\"><code class=\"lang-javascript\"><span class=\"pln\">module</span><span class=\"pun\">.</span><span class=\"pln\">exports </span><span class=\"pun\">=</span><span class=\"pln\"> </span><span class=\"pun\">{</span></code></li><li class=\"L1\"><code class=\"lang-javascript\"><span class=\"pln\">    devServer</span><span class=\"pun\">:</span><span class=\"pln\"> </span><span class=\"pun\">{</span></code></li><li class=\"L2\"><code class=\"lang-javascript\"><span class=\"pln\">        proxy</span><span class=\"pun\">:</span><span class=\"pln\"> </span><span class=\"pun\">{</span></code></li><li class=\"L3\"><code class=\"lang-javascript\"><span class=\"pln\">            </span><span class=\"com\">//配置跨域</span></code></li><li class=\"L4\"><code class=\"lang-javascript\"><span class=\"pln\">            </span><span class=\"str\">\'/api\'</span><span class=\"pun\">:</span><span class=\"pln\"> </span><span class=\"pun\">{</span></code></li><li class=\"L5\"><code class=\"lang-javascript\"><span class=\"pln\">                target</span><span class=\"pun\">:</span><span class=\"pln\"> </span><span class=\"str\">\"http://api.org\"</span><span class=\"pun\">,</span></code></li><li class=\"L6\"><code class=\"lang-javascript\"><span class=\"pln\">                ws</span><span class=\"pun\">:</span><span class=\"pln\"> </span><span class=\"kwd\">true</span><span class=\"pun\">,</span></code></li><li class=\"L7\"><code class=\"lang-javascript\"><span class=\"pln\">                changOrigin</span><span class=\"pun\">:</span><span class=\"pln\"> </span><span class=\"kwd\">true</span><span class=\"pun\">,</span></code></li><li class=\"L8\"><code class=\"lang-javascript\"><span class=\"pln\">                pathRewrite</span><span class=\"pun\">:</span><span class=\"pln\"> </span><span class=\"pun\">{</span></code></li><li class=\"L9\"><code class=\"lang-javascript\"><span class=\"pln\">                    </span><span class=\"str\">\'^/api\'</span><span class=\"pun\">:</span><span class=\"pln\"> </span><span class=\"str\">\'/\'</span></code></li><li class=\"L0\"><code class=\"lang-javascript\"><span class=\"pln\">                </span><span class=\"pun\">}</span></code></li><li class=\"L1\"><code class=\"lang-javascript\"><span class=\"pln\">            </span><span class=\"pun\">}</span></code></li><li class=\"L2\"><code class=\"lang-javascript\"><span class=\"pln\">        </span><span class=\"pun\">}</span></code></li><li class=\"L3\"><code class=\"lang-javascript\"><span class=\"pln\">    </span><span class=\"pun\">}</span></code></li><li class=\"L4\"><code class=\"lang-javascript\"><span class=\"pun\">}</span></code></li></ol></pre>\n<p>这里要把请求路径该一下：<br><code>.get(\"api.org/users\")</code><br>改成:<br><code>.get(\"api/users\")</code></p>\n<h3 id=\"h3-u6CE8u610Fu95EEu9898\"><a name=\"注意问题\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>注意问题</h3><p>如果你是用tp写的接口<br>把    <code>app.config里的　app_debug = true 改为　false</code><br><strong>linux下</strong><br><em>runtime 权限改为777</em></p>\n<pre class=\"prettyprint linenums prettyprinted\" style=\"\"><ol class=\"linenums\"><li class=\"L0\"><code class=\"lang-bash\"><span class=\"pln\">$ sudo chmod </span><span class=\"lit\">777</span><span class=\"pln\"> runtime</span></code></li></ol></pre>\n<p><strong>否则请求的时候会有服务器500错误</strong></p>\n<hr>\n<h2 id=\"h2-u53C2u8003u6587u732E\"><a name=\"参考文献\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>参考文献</h2><p><strong><a href=\"https://cli.vuejs.org/zh/config/#devserver-proxy\">https://cli.vuejs.org/zh/config/#devserver-proxy</a></strong><br><strong><a href=\"https://github.com/chimurai/http-proxy-middleware#proxycontext-config\">https://github.com/chimurai/http-proxy-middleware#proxycontext-config</a></strong></p>\n', 1, '最近学了Vue 想写点demo练习一下 就用tp写了个接口', '/static/upload/20210104/cac88e0db4d4feea676d2a0890ade738.jpg', 2, 41, NULL, 0, 1, 0, 'Vue,前端', 1609730194, 1610000544, NULL);
INSERT INTO `t_article` VALUES (8, '敲开PHP函数的大门 ', '函数结构\n构成部分:\n关键字 function\n函数名 functionName\n参数列表 $a,$b\n函数体 {}\n返回值 return $a,$b\n例：\n\nfunction functionName($a,$b)\n{\n    return $a+$b;\n}\n函数参数\n主要有以下几类：\n\n形参&实参\n默认值\n强类型参数\n可变数量的参数列表\n值传递&引用传递\n形参&实参\n形参:\n\nfunction sum($x,$y /*形参*/){\n    $tmp = $x + $y\n    return $tmp\n}\n实参\n\nfunction sum($x,$y){\n    $tmp = $x + $y\n    return $tmp\n}\n$x = 1;\n$y = 2;\nsum($x,$y /*实参*/)\n注意:\n\n调用函数时，形参和实参的关系是一对一的,实参按顺序赋值给形参，实参可以是一个表达式，传递实参大于形参数量不会报错，反之小于形参会报错。\n\n默认值\n解释:允许函数的参数在调用的时候给定具体的值，如果没有对应的值函数可以是使用默认值作为变量，反之如果有则会替换默认值,缺省值应该放在右边，默认值必须是常量\n\nfunction greet_to_someone($name,$is_formal = fales){\n    if ($is_formal) {\n        return \"Hello\" . $name;\n    } else {\n        return \"Hi\" . $name;\n    }\n}\n强类型参数\n定义：为函数列表中的参数制定类型，如果传入的数据类型不匹配，则抛出 TypeError 异常\n\n支持类型:\n\nclass/interface name PHP5.0.0\narray PHP5.1.0\ncallable PHP5.4.0\nbool,float,int,string PHP7.0.0\nfunction sum (int $x, int $y){\n    return $x+$y;\n}\n可变数量的参数列表\ncode/args_list.php\n\n实现方式 1:\n\n在 php5.5 及更早的版本中，使用 func_num_args(), func_get_arg(), func_get_args()函数实现。\n\nfunc_num_args()\n返回用户传递参数的数量\nfunc_get_arg()\n返回参数列表中的某一项(是一个偏移量从 0 开始如果指定偏移量大于传递偏移量会有警告)\nfunc_get_args()\n返回包含函数参数列表的数组\nfunction get_sum()\n{\n    $args_num = func_num_args();\n    $sum = 0;\n    if ($args_num == 0) {\n        return 0;\n    } else {\n        for ($i = 0; $i < $args_num; $i++) {\n            $sum += func_get_arg($i);\n        }\n        return $sum;\n    }\n}\n\necho get_sum(1,2,4);\n实现方式 2:\n\n在 php5.6 及以上的版本中，可以使用…语法实现,返回一个数组。\n\nfunction get_sum(...$args)\n{\n    if ($args == []) {\n        return 0;\n    } else {\n        for ($i = 0; $i < count($args); $i++) {\n            $sum += $args[$i];\n        }\n        return $sum;\n    }\n}\n\necho get_sum(1, 2, 3);\n值传递&引用传递\ncode/args_pass.php\n\n1. 值传递:\n\nfunction factorial($num)\n{\n    $num = $num * 2;\n    return $num;\n}\n\n$num = 3;\necho factorial($num);\necho $num;\n// 此时我们发现函数并不会影响实际的变量$num\n}\n1. 引用传递:\n\n\nfunction swap(&$x, &$y)\n{\n    $tmp = $x;\n    $x = $y;\n    $y = $tmp;\n}\n\n\n$a = 1;$b = 2;\nswap($a, $b);\n\necho \"\\n\" . $a, $b;\n// 此时我们发现变量$a,$b值发生了交换\n复杂函数\n可变函数\n定义:\n\n变量名后有圆括号,PHP 将寻找与变量值相同名的函数,并且尝试执行它\ncode/var_function.php\n嵌套函数\ncode/out_in.php\n特点:PHP 嵌套函数有一些特别之处。最特别的是，当外部函数被调用时，内部函数就会自动进入全局作用域中，成为新的定义函数\n\n递归函数\n定义：\n\ncode/recursiwe.php\n\n函数在它的函数体调用自身，这种函数称为递归函数\n作用:\n\n分解问题，调用自身\n匿名函数\ncode/Anonymous_functions.php\n\n定义:\n\n匿名函数(Anonymous functions) ，也叫闭包(closures),允许 临时创建一个没有指定函数名的函数，最经常作为回调函数(callback)参数的值\n使用\n\n闭包函数可以作为变量的值来使用。\n代码复用\n如果你想在一个文件里面引入其他文件的函数\n可以使用一下几种方法\n\ninclude()\nrequire()\ninclude 和 require 语句是相同的，除了错误处理方面：\n\nrequire 会生成致命错误（E_COMPILE_ERROR）并停止脚本\ninclude 只生成警告（E_WARNING），并且脚本会继续\n\ninclude_once()\nrequire_once()\n表示文件只能引入一次\n\n引入文件夹\n\nset_include_path();\nget_include_path();\n', '<p>函数结构<br>构成部分:<br>关键字 function<br>函数名 functionName<br>参数列表 $a,$b<br>函数体 {}<br>返回值 return $a,$b<br>例：</p>\n<p>function functionName($a,$b)<br>{<br>    return $a+$b;<br>}<br>函数参数<br>主要有以下几类：</p>\n<p>形参&amp;实参<br>默认值<br>强类型参数<br>可变数量的参数列表<br>值传递&amp;引用传递<br>形参&amp;实参<br>形参:</p>\n<p>function sum($x,$y /<em>形参</em>/){<br>    $tmp = $x + $y<br>    return $tmp<br>}<br>实参</p>\n<p>function sum($x,$y){<br>    $tmp = $x + $y<br>    return $tmp<br>}<br>$x = 1;<br>$y = 2;<br>sum($x,$y /<em>实参</em>/)<br>注意:</p>\n<p>调用函数时，形参和实参的关系是一对一的,实参按顺序赋值给形参，实参可以是一个表达式，传递实参大于形参数量不会报错，反之小于形参会报错。</p>\n<p>默认值<br>解释:允许函数的参数在调用的时候给定具体的值，如果没有对应的值函数可以是使用默认值作为变量，反之如果有则会替换默认值,缺省值应该放在右边，默认值必须是常量</p>\n<p>function greet_to_someone($name,$is_formal = fales){<br>    if ($is_formal) {<br>        return “Hello” . $name;<br>    } else {<br>        return “Hi” . $name;<br>    }<br>}<br>强类型参数<br>定义：为函数列表中的参数制定类型，如果传入的数据类型不匹配，则抛出 TypeError 异常</p>\n<p>支持类型:</p>\n<p>class/interface name PHP5.0.0<br>array PHP5.1.0<br>callable PHP5.4.0<br>bool,float,int,string PHP7.0.0<br>function sum (int $x, int $y){<br>    return $x+$y;<br>}<br>可变数量的参数列表<br>code/args_list.php</p>\n<p>实现方式 1:</p>\n<p>在 php5.5 及更早的版本中，使用 func_num_args(), func_get_arg(), func_get_args()函数实现。</p>\n<p>func_num_args()<br>返回用户传递参数的数量<br>func_get_arg()<br>返回参数列表中的某一项(是一个偏移量从 0 开始如果指定偏移量大于传递偏移量会有警告)<br>func_get_args()<br>返回包含函数参数列表的数组<br>function get_sum()<br>{<br>    $args_num = func_num_args();<br>    $sum = 0;<br>    if ($args_num == 0) {<br>        return 0;<br>    } else {<br>        for ($i = 0; $i &lt; $args_num; $i++) {<br>            $sum += func_get_arg($i);<br>        }<br>        return $sum;<br>    }<br>}</p>\n<p>echo get_sum(1,2,4);<br>实现方式 2:</p>\n<p>在 php5.6 及以上的版本中，可以使用…语法实现,返回一个数组。</p>\n<p>function get_sum(…$args)<br>{<br>    if ($args == []) {<br>        return 0;<br>    } else {<br>        for ($i = 0; $i &lt; count($args); $i++) {<br>            $sum += $args[$i];<br>        }<br>        return $sum;<br>    }<br>}</p>\n<p>echo get_sum(1, 2, 3);<br>值传递&amp;引用传递<br>code/args_pass.php</p>\n<ol>\n<li>值传递:</li></ol>\n<p>function factorial($num)<br>{<br>    $num = $num * 2;<br>    return $num;<br>}</p>\n<p>$num = 3;<br>echo factorial($num);<br>echo $num;<br>// 此时我们发现函数并不会影响实际的变量$num<br>}</p>\n<ol>\n<li>引用传递:</li></ol>\n<p>function swap(&amp;$x, &amp;$y)<br>{<br>    $tmp = $x;<br>    $x = $y;<br>    $y = $tmp;<br>}</p>\n<p>$a = 1;$b = 2;<br>swap($a, $b);</p>\n<p>echo “\\n” . $a, $b;<br>// 此时我们发现变量$a,$b值发生了交换<br>复杂函数<br>可变函数<br>定义:</p>\n<p>变量名后有圆括号,PHP 将寻找与变量值相同名的函数,并且尝试执行它<br>code/var_function.php<br>嵌套函数<br>code/out_in.php<br>特点:PHP 嵌套函数有一些特别之处。最特别的是，当外部函数被调用时，内部函数就会自动进入全局作用域中，成为新的定义函数</p>\n<p>递归函数<br>定义：</p>\n<p>code/recursiwe.php</p>\n<p>函数在它的函数体调用自身，这种函数称为递归函数<br>作用:</p>\n<p>分解问题，调用自身<br>匿名函数<br>code/Anonymous_functions.php</p>\n<p>定义:</p>\n<p>匿名函数(Anonymous functions) ，也叫闭包(closures),允许 临时创建一个没有指定函数名的函数，最经常作为回调函数(callback)参数的值<br>使用</p>\n<p>闭包函数可以作为变量的值来使用。<br>代码复用<br>如果你想在一个文件里面引入其他文件的函数<br>可以使用一下几种方法</p>\n<p>include()<br>require()<br>include 和 require 语句是相同的，除了错误处理方面：</p>\n<p>require 会生成致命错误（E_COMPILE_ERROR）并停止脚本<br>include 只生成警告（E_WARNING），并且脚本会继续</p>\n<p>include_once()<br>require_once()<br>表示文件只能引入一次</p>\n<p>引入文件夹</p>\n<p>set_include_path();<br>get_include_path();</p>\n', 1, '构成部分:\n关键字 function\n函数名 functionName\n参数列表 $a,$b\n函数体 {}\n返回值 return $a,$b\n', '/static/upload/20210104/d67df28508a8f716fa91fca2c0fdfdab.jpg', 2, 1, NULL, 0, 1, 0, '函数,编程', 1609730201, 1610000552, 1);
INSERT INTO `t_article` VALUES (9, '关于我', '\n## 关于此博客\n\n很开心你能发现我，这是我的博客不定期更新文章哦\n\n\n\n关于我博客的一些信息\n\n<img src=\"http://blog.im/static/upload/20210104/e3725d59b6649d14f18ee1bbd6f5011f.jpeg\" style=\"zoom: 25%;\" >\n\n\n\n博客：这套博客源码是我自己用PHP写的，并且我把他开源了！！！\n\n开源地址: https://github.com/muxik/PureBlog\n\n\n\n域名: muxik.top 2020年十二月购入\n\n服务器: 在cloud.cn 买了一年\n\n以后也会继续续费的\n\n所以 不用担心Muxi_k会突然消失哦(๑•̀ㅂ•́)و✧\n\n\n\n## 关于我\n\n\n\n我是一个零零后哦，想不到吧!\n\n生日呢是在7月， 具体就不说啦！~~以后在告诉你们吧!~~\n\n\n\n**爱好的话**\n\n有很多呀，但是好多都东西都没有去做! ~~可能是时间太少~~\n\n\n\n- 编程\n\n  一些经历\n\n  ​		很小的时候就对黑客感兴趣，16年的时候去各种论坛学习~~黑客技术~~ 就是一些脚本的用法，像啊D，明小子，御剑之类的那种\n\n  后来不知什么时候遇到了一个~~道上的~~大佬，他告诉我想学黑客技术要先掌握好编程。\n\n  ​		后来我就走上了编程的道路，把黑客什么的忘的一干二净，我最开始学的是 C ，那个时候去爱奇艺找的教程视频都看的一知半懂，最后干脆放弃了\n\n  后来又去学了Python 终于看的懂一点了，然后我买了本Python的书，开始学了起来，学着书里的小例子，很感兴趣！\n\n  ​		18年年末我搭建了一个网站:~~muxikk.top~~ 当时服务器带着我的数据跑路了，域名也就没在续费了，有点小遗憾吧，但是呢也从搭建网站的过程中学到了很\n\n  多东西，前端后端都有，是一次不错的经历\n\n  ​		再后来社团要做一个招新系统，当时我正在学PHP就帮社团做了一个，也就是CURD和前端的一些页面当时用了一些前端框架，bootstrap, layui,很快就写好了。算是第一次帮别人写东西吧 \n\n  \n\n- GNU/Linux\n\n  ​		很早之前就接触过一些发行版，kali,ubuntu, 开始正式接触的时候是有一次使用php做开发的时候，遇到了一些环境依赖的问题是Windows的造成的结果，我就果断换了个Fedora问题立马就解决了，当时真的是太开心了,再后来大概用了一段时间的Windows因为当时要用一些Adobe的东西当时就装了双系统，后来windows又出问题了，就把两个系统都格了，看一个B站up@TheCW的视频就换了Manjaro,使用期间又换了Deepin,Ubuntu,opensuse,最后定居在了Arch 现在已经用了快一年半了，体验很不错\n\n  ​		很喜欢Linux的哲学以及高度自由的操作和一切皆文件的哲学\n\n  \n\n- 画画\n\n  ​		在很久很久一前曾跟随杨世相老师学过一段时间的水墨画，虽然\n\n  时间很短，但是也学到了很多东西！在这里感谢杨老师。\n\n  \n\n- 吉他\n\n  一直很喜欢宋胖子就跟室友学了吉他，和一点乐理，弹的很烂，就不多提了', '<h2 id=\"h2-u5173u4E8Eu6B64u535Au5BA2\"><a name=\"关于此博客\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>关于此博客</h2><p>很开心你能发现我，这是我的博客不定期更新文章哦</p>\n<p>关于我博客的一些信息</p>\n<p>&lt;img src=\"http://blog.im/static/upload/20210104/e3725d59b6649d14f18ee1bbd6f5011f.jpeg\" style=\"zoom: 25%;\" &gt;</p>\n<p>博客：这套博客源码是我自己用PHP写的，并且我把他开源了！！！</p>\n<p>开源地址: <a href=\"https://github.com/muxik/PureBlog\">https://github.com/muxik/PureBlog</a></p>\n<p>域名: muxik.top 2020年十二月购入</p>\n<p>服务器: 在cloud.cn 买了一年</p>\n<p>以后也会继续续费的</p>\n<p>所以 不用担心Muxi_k会突然消失哦(๑•̀ㅂ•́)و✧</p>\n<h2 id=\"h2-u5173u4E8Eu6211\"><a name=\"关于我\" class=\"reference-link\"></a><span class=\"header-link octicon octicon-link\"></span>关于我</h2><p>我是一个零零后哦，想不到吧!</p>\n<p>生日呢是在7月， 具体就不说啦！<del>以后在告诉你们吧!</del></p>\n<p><strong>爱好的话</strong></p>\n<p>有很多呀，但是好多都东西都没有去做! <del>可能是时间太少</del></p>\n<ul>\n<li><p>编程</p>\n<p>一些经历</p>\n<p>​        很小的时候就对黑客感兴趣，16年的时候去各种论坛学习<del>黑客技术</del> 就是一些脚本的用法，像啊D，明小子，御剑之类的那种</p>\n<p>后来不知什么时候遇到了一个<del>道上的</del>大佬，他告诉我想学黑客技术要先掌握好编程。</p>\n<p>​        后来我就走上了编程的道路，把黑客什么的忘的一干二净，我最开始学的是 C ，那个时候去爱奇艺找的教程视频都看的一知半懂，最后干脆放弃了</p>\n<p>后来又去学了Python 终于看的懂一点了，然后我买了本Python的书，开始学了起来，学着书里的小例子，很感兴趣！</p>\n<p>​        18年年末我搭建了一个网站:<del>muxikk.top</del> 当时服务器带着我的数据跑路了，域名也就没在续费了，有点小遗憾吧，但是呢也从搭建网站的过程中学到了很</p>\n<p>多东西，前端后端都有，是一次不错的经历</p>\n<p>​        再后来社团要做一个招新系统，当时我正在学PHP就帮社团做了一个，也就是CURD和前端的一些页面当时用了一些前端框架，bootstrap, layui,很快就写好了。算是第一次帮别人写东西吧 </p>\n</li></ul>\n<ul>\n<li><p>GNU/Linux</p>\n<p>​        很早之前就接触过一些发行版，kali,ubuntu, 开始正式接触的时候是有一次使用php做开发的时候，遇到了一些环境依赖的问题是Windows的造成的结果，我就果断换了个Fedora问题立马就解决了，当时真的是太开心了,再后来大概用了一段时间的Windows因为当时要用一些Adobe的东西当时就装了双系统，后来windows又出问题了，就把两个系统都格了，看一个B站up<a href=\"https://github.com/TheCW\" title=\"@TheCW\" class=\"at-link\"></a><a href=\"https://github.com/TheCW\" title=\"@TheCW\" class=\"at-link\">@TheCW</a>的视频就换了Manjaro,使用期间又换了Deepin,Ubuntu,opensuse,最后定居在了Arch 现在已经用了快一年半了，体验很不错</p>\n<p>​        很喜欢Linux的哲学以及高度自由的操作和一切皆文件的哲学</p>\n</li></ul>\n<ul>\n<li><p>画画</p>\n<p>​        在很久很久一前曾跟随杨世相老师学过一段时间的水墨画，虽然</p>\n<p>时间很短，但是也学到了很多东西！在这里感谢杨老师。</p>\n</li></ul>\n<ul>\n<li><p>吉他</p>\n<p>一直很喜欢宋胖子就跟室友学了吉他，和一点乐理，弹的很烂，就不多提了</p>\n</li></ul>\n', 1, '一些关于，此博客的说明和我的小经历', '/static/upload/20210105/6dcc49fa81affe5ce7e8eb8cf74cef75.jpg', 2, 0, NULL, 0, 1, 0, '关于', 1609844307, 1609844307, 1609857555);
COMMIT;

-- ----------------------------
-- Table structure for t_category
-- ----------------------------
DROP TABLE IF EXISTS `t_category`;
CREATE TABLE `t_category` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `link` varchar(255) DEFAULT NULL,
  `type` tinyint(1) DEFAULT '0',
  `sort` int DEFAULT '0',
  `pid` int NOT NULL,
  `state` tinyint(1) NOT NULL DEFAULT '1',
  `create_time` int DEFAULT NULL,
  `update_time` int DEFAULT NULL,
  `delete_time` int DEFAULT NULL,
  `page` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_category
-- ----------------------------
BEGIN;
INSERT INTO `t_category` VALUES (1, '首页', '/', 1, 1, 0, 1, 1609167190, 1, NULL, NULL);
INSERT INTO `t_category` VALUES (2, '开发', NULL, 0, 2, 0, 1, 1609167190, 1, NULL, NULL);
INSERT INTO `t_category` VALUES (3, '生活', NULL, 0, 2, 0, 1, 1609167190, 1, NULL, NULL);
INSERT INTO `t_category` VALUES (4, '友人', '/link.html', 1, 3, 0, 1, 1609167190, 1609581220, NULL, NULL);
INSERT INTO `t_category` VALUES (5, 'PHP', NULL, 0, 2, 2, 1, 1609167190, 1610108253, 1610108253, NULL);
INSERT INTO `t_category` VALUES (6, 'JavaScript', NULL, 0, 1, 2, 1, 1609167190, 1610108251, 1610108251, NULL);
INSERT INTO `t_category` VALUES (7, '琐事\n', NULL, 0, 2, 3, 1, 1609167190, 1, NULL, NULL);
INSERT INTO `t_category` VALUES (17, '关于', '/about.html', 1, 3, 0, 1, 1609167190, 1, NULL, NULL);
INSERT INTO `t_category` VALUES (19, '碎碎念', '/page/time_line.html', 1, 3, 3, 1, 1609167190, 1609167190, NULL, NULL);
INSERT INTO `t_category` VALUES (20, 'Linux', NULL, 0, 0, 2, 1, 1609675722, 1609675858, 1609675858, NULL);
INSERT INTO `t_category` VALUES (21, 'Linux', NULL, 0, 0, 0, 1, 1609675867, 1609675895, 1609675895, 'linux');
INSERT INTO `t_category` VALUES (22, 'GNU/Linux', NULL, 0, 0, 2, 1, 1609675920, 1610108248, 1610108248, 'linux');
INSERT INTO `t_category` VALUES (23, 'PHP', NULL, 0, 0, 2, 1, 1610108286, 1610108286, NULL, 'php');
INSERT INTO `t_category` VALUES (24, '数据库', NULL, 0, 0, 2, 1, 1610108313, 1610108313, NULL, 'database');
COMMIT;

-- ----------------------------
-- Table structure for t_comment
-- ----------------------------
DROP TABLE IF EXISTS `t_comment`;
CREATE TABLE `t_comment` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `page_id` int DEFAULT NULL,
  `article_id` int DEFAULT NULL,
  `pid` int DEFAULT '0',
  `os` varchar(255) DEFAULT NULL,
  `browser` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `content` varchar(255) DEFAULT NULL,
  `site` varchar(255) DEFAULT NULL,
  `nickname` varchar(255) DEFAULT NULL,
  `create_time` int DEFAULT NULL,
  `update_time` int DEFAULT NULL,
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=72 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_comment
-- ----------------------------
BEGIN;
INSERT INTO `t_comment` VALUES (32, 1, NULL, 0, 'Linux', 'Chrome', 'a', 'a', 'a', 'a', 1609991944, 1609991944, NULL);
INSERT INTO `t_comment` VALUES (33, 1, NULL, 32, 'Linux', 'Chrome', 'dasd', 'dasd', 'asdas', 'asdas', 1609991956, 1609991956, NULL);
INSERT INTO `t_comment` VALUES (34, 2, NULL, 0, 'Linux', 'Chrome', 'sdf', 'sdfs', 'sdf', 'sdf', 1609992137, 1609992137, NULL);
INSERT INTO `t_comment` VALUES (35, 2, NULL, 0, 'Linux', 'Chrome', 'sdf', 'dsfsd', 'sdf', 'sdf', 1609992156, 1609992156, NULL);
INSERT INTO `t_comment` VALUES (36, 2, NULL, 34, 'Linux', 'Chrome', 'asd', 'dasd', 'sdas', 'sad', 1609992194, 1609992194, NULL);
INSERT INTO `t_comment` VALUES (37, NULL, 8, 0, 'Linux', 'Chrome', 'asda', 'sdasd', 'das', 'sad', 1609999727, 1610000552, 1610000552);
INSERT INTO `t_comment` VALUES (38, NULL, 7, 0, 'Linux', 'Chrome', 'lqjxm666@163.com', 'sdfsdfsdf', 'muxik.top', 'Muxi_k', 1610007467, 1610007467, NULL);
INSERT INTO `t_comment` VALUES (39, NULL, 7, 38, 'Linux', 'Chrome', 'kjskdljfk', 'kdjlfkl', 'kdjkgsdl', 'sflsldfl;', 1610007479, 1610007479, NULL);
INSERT INTO `t_comment` VALUES (40, 1, NULL, 0, 'Linux', 'Chrome', 'asd', 'asda', 'asd', 'asd', 1610007685, 1610007685, NULL);
INSERT INTO `t_comment` VALUES (41, 1, NULL, 32, 'Linux', 'Chrome', 'lqjxm666@163.com', 'asd', '', 'asdasd', 1610007694, 1610007694, NULL);
INSERT INTO `t_comment` VALUES (42, 2, NULL, 0, 'Linux', 'Chrome', 'lqjxm666@163.com', 'aasfs', 'http://www.muxik.top', 'Muxi_k', 1610026347, 1610026347, NULL);
INSERT INTO `t_comment` VALUES (43, 1, NULL, 0, 'Linux', 'Chrome', 'sdf', 'dfsd', 'http://sds', 'sdf', 1610095903, 1610095903, NULL);
INSERT INTO `t_comment` VALUES (44, 1, NULL, 0, 'Linux', 'Chrome', 'fffff', 'ffff', 'http://fffff', 'ffff', 1610095912, 1610095912, NULL);
INSERT INTO `t_comment` VALUES (45, 2, NULL, 0, 'Linux', 'Chrome', 'asd', 'asdasd', 'http://', 'asd', 1610096010, 1610096010, NULL);
INSERT INTO `t_comment` VALUES (46, NULL, 7, 0, 'Linux', 'Chrome', 'sdf', 'fsdf', 'http://sdfsd', 'sdf', 1610096206, 1610096206, NULL);
INSERT INTO `t_comment` VALUES (47, NULL, 7, 46, 'Linux', 'Chrome', 'dfsdf', 'sdfsdf', 'sdfhttp://', 'dsfs', 1610096215, 1610096215, NULL);
INSERT INTO `t_comment` VALUES (48, 2, NULL, 42, 'Linux', 'Chrome', 'asd', 'asda', 'asdhttp://', 'asdasd', 1610096328, 1610096328, NULL);
INSERT INTO `t_comment` VALUES (49, 2, NULL, 42, 'Linux', 'Chrome', 'fsd', 'sdfsd', 'http://sdf', 'dsfsd', 1610096339, 1610096339, NULL);
INSERT INTO `t_comment` VALUES (50, 3, NULL, 0, 'Linux', 'Chrome', 'test', 'kddfkj', 'http://', 'test', 1610205076, 1610205076, NULL);
INSERT INTO `t_comment` VALUES (51, 3, NULL, 0, 'Linux', 'Chrome', 'sdfs', 'dfsdf', 'http://', 'sdf', 1610205227, 1610205227, NULL);
INSERT INTO `t_comment` VALUES (52, 1, NULL, 43, 'Linux', 'Chrome', 's d f', 'sdf', 's d fhttp://', 's d f', 1610287445, 1610287445, NULL);
INSERT INTO `t_comment` VALUES (53, 2, NULL, 42, 'Linux', 'Chrome', 'a', 'a', 'http://a', 'a', 1610331305, 1610331305, NULL);
INSERT INTO `t_comment` VALUES (54, NULL, 7, 46, 'Linux', 'Chrome', 'dfs', 'sdf', 'http://', 'sdf', 1610331371, 1610331371, NULL);
INSERT INTO `t_comment` VALUES (55, NULL, 7, 0, 'Linux', 'Chrome', 'asd', 'asdasd', 'http://', 'asf', 1610331379, 1610331379, NULL);
INSERT INTO `t_comment` VALUES (56, NULL, 7, 38, 'Linux', 'Chrome', 'sdf', 'sdfsd', 'http://', 'sdf', 1610331524, 1610331524, NULL);
INSERT INTO `t_comment` VALUES (57, NULL, 7, 0, 'Linux', 'Chrome', 'asd', 'sadasd', 'http://', 'asdaaa', 1610332259, 1610332259, NULL);
INSERT INTO `t_comment` VALUES (58, 3, NULL, 0, 'Linux', 'Chrome', 'a', 'a', 'http://', 'a', 1610333023, 1610333023, NULL);
INSERT INTO `t_comment` VALUES (59, 1, NULL, 0, 'Linux', 'Chrome', 'fsdf', 'sdfsdf', 'http://', 'sdfsd', 1610333260, 1610333260, NULL);
INSERT INTO `t_comment` VALUES (60, 1, NULL, 0, 'Linux', 'Chrome', 'aaaaaa', 'aaaa', 'http://aaaaa', 'aaaa', 1610333300, 1610333300, NULL);
INSERT INTO `t_comment` VALUES (61, 2, NULL, 45, 'Linux', 'Chrome', 'asd', 'dsadsa', 'http://', 'asd', 1610333557, 1610333557, NULL);
INSERT INTO `t_comment` VALUES (62, 3, NULL, 0, 'Linux', 'Chrome', 'sdfsd', 'fsdf', 'http://', 'sdf', 1610333610, 1610333610, NULL);
INSERT INTO `t_comment` VALUES (63, 3, NULL, 0, 'Linux', 'Chrome', 'sdfsd', 'fsdf', 'http://', 'sdf', 1610333610, 1610333610, NULL);
INSERT INTO `t_comment` VALUES (64, 3, NULL, 0, 'Linux', 'Chrome', 'sdfsd', 'fsdf', 'http://', 'sdf', 1610333610, 1610333610, NULL);
INSERT INTO `t_comment` VALUES (65, 3, NULL, 0, 'Linux', 'Chrome', 'aaaaaaaaaa', 'aaaaaaaaaaaaa', 'http://', 'asdasaaaaaaaaaaa', 1610333634, 1610333634, NULL);
INSERT INTO `t_comment` VALUES (66, 3, NULL, 65, 'Linux', 'Chrome', 'a', 'a', 'ahttp://', 'aa', 1610333644, 1610333644, NULL);
INSERT INTO `t_comment` VALUES (67, 3, NULL, 65, 'Linux', 'Chrome', 'c', 'c', 'http://', 'c', 1610333653, 1610333653, NULL);
INSERT INTO `t_comment` VALUES (68, NULL, 7, 0, 'Linux', 'Chrome', 'aaa', 'aaa', 'http://', 'aaaaaaa', 1610333781, 1610333781, NULL);
INSERT INTO `t_comment` VALUES (69, NULL, 7, 60, 'Linux', 'Chrome', 'dd', 'ddddddddddd', 'http://dd', 'dddddd', 1610334444, 1610334444, NULL);
INSERT INTO `t_comment` VALUES (70, NULL, 7, 60, 'Linux', 'Chrome', 'dada', 'asdasdas', 'http://adasd', 'asdas', 1610334457, 1610334457, NULL);
INSERT INTO `t_comment` VALUES (71, NULL, 7, 60, 'Linux', 'Chrome', 'fsd', 'sdfsd', 'http://sdf', 'dd', 1610334502, 1610334502, NULL);
COMMIT;

-- ----------------------------
-- Table structure for t_link
-- ----------------------------
DROP TABLE IF EXISTS `t_link`;
CREATE TABLE `t_link` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `type` tinyint(1) DEFAULT '1',
  `logo` varchar(255) DEFAULT NULL,
  `link` varchar(255) DEFAULT NULL,
  `state` tinyint(1) DEFAULT '1',
  `create_time` int DEFAULT NULL,
  `update_time` int DEFAULT NULL,
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_link
-- ----------------------------
BEGIN;
INSERT INTO `t_link` VALUES (1, '千千', '继续踏上旅途，在没有你的春天……', 0, 'https://www.dreamwings.cn/wp-content/uploads/2018/06/806e52a2e2b9ff4bd2c23140df75cc1f.jpeg', 'https://www.dreamwings.cn', 1, 1609341326, 1609341326, NULL);
INSERT INTO `t_link` VALUES (2, '顾思维', '啊啊啊', 0, 'https://19930.vip/logo.png', 'https://19930.vip/', 0, 1609379066, 1609379066, NULL);
INSERT INTO `t_link` VALUES (3, '菜鸟教程', '学习使我快乐', 0, 'https://static.runoob.com/images/favicon.ico', 'https://www.runoob.com/', 0, 1609380926, 1609380926, NULL);
INSERT INTO `t_link` VALUES (4, 'bilibili', 'Muxi_k的bilibili空间', 1, '/static/upload/20210105/cfa7c6fbd319cf48d7be2c6ad343a01b.png', 'https://space.bilibili.com/295239068', 1, 1609825266, 1609825266, NULL);
INSERT INTO `t_link` VALUES (5, 'Facebook', 'Muxi_k的FaceBook地址', 1, '/static/upload/20210105/cabf0b83c62a2bac20a56c6e2cbca1f3.png', 'https://www.facebook.com/profile.php?id=100028304134989', 1, 1609825724, 1609825724, NULL);
INSERT INTO `t_link` VALUES (6, '知乎', 'Muxi_k 的知乎地址', 1, '/static/upload/20210105/a648fa239f3e7671843697539a7af40d.png', 'https://www.zhihu.com/people/skruy', 1, 1609826168, 1609826168, NULL);
INSERT INTO `t_link` VALUES (7, 'Twitter', 'Muxi_K 的Twitter地址', 1, '/static/upload/20210105/62ecebc2c0216c49b0f5fae98f472fa6.png', 'https://twitter.com/muxi_king', 1, 1609826308, 1609826308, NULL);
COMMIT;

-- ----------------------------
-- Table structure for t_shuo
-- ----------------------------
DROP TABLE IF EXISTS `t_shuo`;
CREATE TABLE `t_shuo` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(255) NOT NULL,
  `create_time` int NOT NULL,
  `delete_time` int DEFAULT NULL,
  `year` int DEFAULT NULL,
  `state` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_shuo
-- ----------------------------
BEGIN;
INSERT INTO `t_shuo` VALUES (1, 'aaaaaa', 1610216030, NULL, 2021, 1);
INSERT INTO `t_shuo` VALUES (2, 'sdfsdf', 1610216079, NULL, 2021, 1);
INSERT INTO `t_shuo` VALUES (3, 'jjnlll', 1610330857, NULL, 2021, 1);
COMMIT;

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(36) NOT NULL,
  `password` varchar(36) NOT NULL,
  `nickname` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `state` tinyint(1) NOT NULL DEFAULT '1',
  `create_time` int NOT NULL,
  `update_time` int NOT NULL,
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_user
-- ----------------------------
BEGIN;
INSERT INTO `t_user` VALUES (1, 'user1', '21232f297a57a5a743894a0e4a801fc3', 'admin', 'lqjxm666@163.com', 0, 1609480927, 1609480927, 1610097475);
COMMIT;

-- ----------------------------
-- Table structure for t_web
-- ----------------------------
DROP TABLE IF EXISTS `t_web`;
CREATE TABLE `t_web` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `title` varchar(255) DEFAULT NULL,
  `like` int DEFAULT '0',
  `motto` varchar(255) DEFAULT NULL,
  `keyword` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `logo` varchar(255) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `close` tinyint(1) NOT NULL DEFAULT '0',
  `copyright` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `create_time` int NOT NULL,
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_web
-- ----------------------------
BEGIN;
INSERT INTO `t_web` VALUES (1, 'Muxi_k\'s Blog', 'Muxi_k\'s Blog', 7, '路虽远，行则必达 事虽难，做则必成', 'Muxi_k,muxik,Muxi_k\' Blog', 'Muxi_k 的小站', '/static/upload/20210104/e3725d59b6649d14f18ee1bbd6f5011f.jpeg', 'lqjxm666@163.com', 0, 'Copyright ©2020 PureBlog v2.0 All Rights Reserved', 1, NULL);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
