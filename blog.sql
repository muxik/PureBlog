/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80021
 Source Host           : localhost:3306
 Source Schema         : blog

 Target Server Type    : MySQL
 Target Server Version : 80021
 File Encoding         : 65001

 Date: 05/11/2020 22:31:09
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for blog_admin
-- ----------------------------
DROP TABLE IF EXISTS `blog_admin`;
CREATE TABLE `blog_admin` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(16) NOT NULL,
  `password` varchar(36) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `super` int DEFAULT NULL,
  `status` int DEFAULT '1' COMMENT '1启用0禁用',
  `create_time` int NOT NULL,
  `update_time` int DEFAULT NULL,
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of blog_admin
-- ----------------------------
BEGIN;
INSERT INTO `blog_admin` VALUES (1, 'admin', '21232f297a57a5a743894a0e4a801fc3', 1, 1, 1604037851, 1604562429, NULL);
INSERT INTO `blog_admin` VALUES (2, 'test1', 'e10adc3949ba59abbe56e057f20f883e', 0, 1, 1604117061, 1604201076, NULL);
INSERT INTO `blog_admin` VALUES (3, 'test2', '8ad8757baa8564dc136c1e07507f4a98', 0, 1, 1604117061, 1604141955, NULL);
INSERT INTO `blog_admin` VALUES (4, 'test3', '8ad8757baa8564dc136c1e07507f4a98', 1, 1, 1604117153, 1604562430, NULL);
INSERT INTO `blog_admin` VALUES (5, 'test4', '86985e105f79b95d6bc918fb45ec7727', 1, 0, 1604117202, 1604562537, NULL);
INSERT INTO `blog_admin` VALUES (6, 'test6', '4cfad7076129962ee70c36839a1e3e15', 0, 1, 1604117218, 1604150711, NULL);
INSERT INTO `blog_admin` VALUES (7, 'test', '8ad8757baa8564dc136c1e07507f4a98', 0, 1, 1604129783, 1604129783, NULL);
INSERT INTO `blog_admin` VALUES (8, 'test34', '8ad8757baa8564dc136c1e07507f4a98', 0, 1, 1604129792, 1604129792, NULL);
INSERT INTO `blog_admin` VALUES (9, 'test5', '8ad8757baa8564dc136c1e07507f4a98', 0, 1, 1604129801, 1604129801, NULL);
INSERT INTO `blog_admin` VALUES (10, 'abcd', 'e2fc714c4727ee9395f324cd2e7f331f', 0, 1, 1604145115, 1604150699, NULL);
INSERT INTO `blog_admin` VALUES (11, 'muxadmin', '21232f297a57a5a743894a0e4a801fc3', 1, 1, 1604562239, 1604562486, NULL);
COMMIT;

-- ----------------------------
-- Table structure for blog_article
-- ----------------------------
DROP TABLE IF EXISTS `blog_article`;
CREATE TABLE `blog_article` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `desc` text,
  `top` int NOT NULL DEFAULT '0',
  `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `state` int NOT NULL DEFAULT '1',
  `u_id` int NOT NULL,
  `tag` varchar(255) NOT NULL,
  `category_id` int NOT NULL,
  `read` int DEFAULT '0',
  `pic` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `create_time` int DEFAULT NULL,
  `update_time` int DEFAULT NULL,
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of blog_article
-- ----------------------------
BEGIN;
INSERT INTO `blog_article` VALUES (1, '我的TODO列表', '我的TODO列表', 1, '# PureBlog TODO\n\n\n\n#### 基本模板\n\n- [x] 博客模板\n- [x] 后台模板\n\n---\n\n#### 登录模块\n\n- [x] 后台登录\n- [x] 后台主页\n- [x] 后台退出登录\n\n---\n\n#### 管理员模块\n\n- [x] 管理员列表\n- [x] 管理员权限修饰\n- [x] 管理员添加\n- [x] 管理员修改\n- [x] 管理员分页\n- [x] 管理员搜索\n- [x] 管理员删除\n- [x] 管理员状态修改\n\n---\n\n#### 栏目模块\n\n- [x] 栏目列表\n- [x] 栏目添加\n- [x] 栏目修改\n- [x] 栏目删除\n- [x] 栏目分页\n- [x] 栏目排序\n- [x] 栏目状态修改\n\n---\n\n#### 文章模块\n\n- [x] 文章列表\n- [ ] 文章添加\n- [ ] 文章修改\n- [ ] 文章删除\n- [ ] 文章分页\n- [ ] TODO\n\n                                ', 1, 1, 'todo', 1, 99, '/static/upload/20201104/860c72de21f4cc40daf5c537af88ae9c.jpg', 1, NULL, NULL);
INSERT INTO `blog_article` VALUES (2, '忘穿流年', '无尽的蔚蓝海域，连接起银河两端的思念，我淡漠的忘穿流年，只为你的承诺，我们的七年，我等\n你……', 1, '　无尽的蔚蓝海域，连接起银河两端的思念，我淡漠的忘穿流年，只为你的承诺，我们的七年，我等\n你……\n　　我们都这般年少，轻狂，扬扬洒洒着青春，我们端不起一杯清明，却放不下一盘昼夜，我抬头凝\n望，远方的星辰，是你么？是我们搁浅了的曾经吗？你会回来么…\n　　你曾说过，美是雪月风花，何惧刹那芳华，我不懂，轻轻地望着你，你温柔的笑，无尽的溺爱，醉\n倒了千年的尘封，我们，一定，上世相逢！\n　　还记得吗？我们许诺过的桃花源，我们答应过的要建立那样一个地方，自此与世无争，寡欲今夕，\n你知道么，我也梦到了，好美…\n　　记得我们清风过，卷云离的日子吗？那时，好不快活，我们都说，人生得你，夫已无求。\n　　而今，你在哪里？你找不到回来的路了吗？你不记得我们的一切了吗?是上帝怎样的戏弄，是时光\n怎样的无情，是地母怎样的孤独，是深海怎样的寒冷，是南极怎样的麻木，是我们怎样的不舍，让你舍\n得遗忘了这一切？\n　　我不知道，我只知道，我的眼泪，早已翻涌成灾，我的理智，早已失去了频度。没有你的世界，怎\n样的缤纷都是灰色的。我不知道你会不会知道，这里的一字一句都是我用一滴滴思念凝成的字眼，我不\n猜你喜欢\n抒情散文 名家抒情散文\n思念句子 写景抒情散文\n想念的话 服务理念口号\n感念师恩作文 怀念的反义词\n思念亲人的古诗 怀念母亲教案\n理想与信念演讲稿 怀念母亲教案\n怀念父亲的诗 秋天的怀念教案\n秋天的怀念教学反思 对已逝亲人思念的句\n自\n2020/11/4 念抒情散文\nwww.ruiwen.com/sanwen/2576994.html 3/12\n知道你会不会知道，在这样一处天堂，我期待着你的归来，我不知道，多少年的守候有没有用，但我等\n着你，就像海岸一样，无论我将见证多少次离别，我都不会离开。我们七年的承诺，你记得吗？没关\n系，我来守，我守着我们的桃源，我守着我们的梦，我守着我放不下的一切，七年，十年，二十年，我\n都等着…\n　　有人说，这叫执念，是心魔，是梦魇，我淡淡的笑…是执念又怎样呢？没有了执念，我们活着又是\n为了什么呢？这一切都不重要了，如果，是心魔，那我愿与它一起堕入无尽的深渊，因为那里，一定有\n我们的梦。如果，是梦魇，我愿意被吞噬，至少，还有孤独的夜色作伴，此生，也值了。如果你在，一\n定也这么想。\n　　执念也罢，残浮也好，若浮生是梦，那我愿梦浮此生。\n念抒情散文4\n　　最初的梦想——这个世界没有一寸干净的土壤，唯有守住心灵的一方净土。我希望我有一支比张爱\n玲更犀利的笔，给更多挣扎在苦痛边缘的人力量和慰籍。\n　　开始喜欢一个叫几米的人，他可以左右一个人的天地。房子高举在云端，猫在吹笛。我忽然想要有\n这样的童心才能为开满花的世界执笔。悲伤，淡淡。\n　　可怜的花，我虏你到我掌心，破了你的世界，碎了你的梦，我也必好生待你。众生本不平等，奈若\n何，生活再乱又怎么样，花点时间一样可以理清。\n　　你千里迢迢而来，只为告诉我一句:对不起，我要忘了你。我对你说，从渺茫中来，到渺茫中去\n吧，幸福在远方，不在这里。你要幸福，不要回头。\n　　人总是在包容，一个人的时候包容另一个遗憾的自己，两个人的时候包容对方，一群人的时候包容\n每一个人。\n　　也许你知道，我的世界从来是不允许不完美存在的，譬如一件心爱的东西，碎一点点也宁愿让它粉\n身碎骨，不为瓦全的。后来明白，万事万物均不得周全，残缺才是真正的美，完整的都是假的。\n　　午后的阳光暖起来，我蓦然发现再寒冷的冬天，也终会有过去的一天。\n　　傣妹与肯德基，居然只是面对面的距离，却又是只能回望的距离。\n　　有时候听着一首歌，突然就想流泪了。并不是歌词有多么感人，只是一瞬间的旋律，美得漫无边\n际。记得京华烟云里，木兰看到美得东西，也总会流泪的。\n　　荒草没幽径，物已非，人又非，谁家心事空迂回？\n　　如果我们可以活到80岁，那么生命的四分之三能与另一个人一起渡过，将是多么幸福的事情。\n　　又看到了长阳的绿水青山，江山如此多娇。只是不知道划船到对岸还是不是一块钱，反正吃到了五\n毛两个的可乐冰。\n　　走在火车站的路上，看到候车室外遍地的人群，大包小包的行李，满脸的悴容，我在想这就是生命\n啊！无论身在哪里，家都是永恒不变的眷恋。忽然想三十岁的我该是哪副模样，是不是一样在拥挤的车\n站挤火车。\n　　黛玉最终，泪尽而逝。 一切皆宿命。 银屏说，一个人死了，有人在她坟前真心流一滴泪，她就是\n值得的。 我将不再为谁哭泣。你说，不要哭。\n　　“北在哪里？”“北在我怀里”！\n　　看完关云长，有一句经典的台词：关羽对曹操说，因为你们这里没有恨，也没有情。生逢乱世，根\n本就没有真正的道义。\n　　曾以为会一辈子刻骨铭心的东西，竟然在念念不忘中忘却了；曾以为会一辈子惺惺相惜的东西，竟\n然在不知不觉中丢失了。很想知道、还剩下什么……华丽的遇见？路到最后，黯淡散场。\n　　有些男人自以为自己一个人功成名就就是真男人，殊不知真男人既铁骨铮铮，又柔情万丈。\n　　风微凉，静静的享受一杯奶茶的时光，安宁而幸福。记得谁说过岁月静好。我也只想生生欢颜。                ', 1, 1, 'test|html', 2, 2, '/static/upload/20201104/ba570967f03951b2d941b987c6f8e544.jpg', 1604192552, NULL, NULL);
INSERT INTO `blog_article` VALUES (3, 'MuxAdmin-简约美观后台模板', '一个简约美观的通用后台模板，基于mdui layui框开发', 1, '一个简约美观的通用后台模板，基于mdui layui框开发\n\n\n\n### 说明\n\n> 目前作者仅提供html版本iframe版本请自行改造\n\n- html版本\n  - 可以进行快速的模板分离\n  - 快速构建后端项目\n\n### 下载\n\n\n$ git clone git@github.com:muxik/MuxAdmin.git\n\n\n                                                ', 1, 1, '后台模板|MDUI|Layui', 3, 3, '/static/upload/20201104/5bad772e517190fa3604a8b72a9c2f13.png', 1604211524, 1604211540, NULL);
INSERT INTO `blog_article` VALUES (4, 'rtklerkt', NULL, 1, '# Test Arcticle', 1, 1, 'tstst', 4, 0, '/static/upload/20201101/85a1c3949f45d197d1c4df1f79d0ff95.jpg', 1604217858, 1604554092, 1604554092);
INSERT INTO `blog_article` VALUES (5, '春', '盼望着，盼bai望着，东风来了，春天的脚步近了。', 0, '盼望着，盼望着，东风来了，春天的脚步近了。\n\n一切都像刚睡醒的样子，欣欣然张开了眼。山朗润起来了，水涨起来了，太阳的脸红起来了。小草偷偷地从土里钻出来，嫩嫩的，绿绿的。园子里，田野里，瞧去，一大片一大片满是的。坐着，躺着，打两个滚，踢几脚球，赛几趟跑，捉几回迷藏。风轻悄悄的，草软绵绵的。\n\n桃树、杏树、梨树，你不让我，我不让你，都开满了花赶趟儿。红的像火，粉的像霞，白的像雪。花里带着甜味儿；闭了眼，树上仿佛已经满是桃儿、杏儿、梨儿。\n\n花下成千成百的蜜蜂嗡嗡地闹着，大小的蝴蝶飞来飞去。野花遍地是：杂样儿，有名字的，没名字的，散在草丛里，像眼睛，像星星，还眨呀眨的。\n\n“吹面不寒杨柳风”，不错的，像母亲的手抚摸着你。风里带来些新翻的泥土的气息，混着青草味儿，还有各种花的香，都在微微润湿的空气里酝酿。\n\n鸟儿将巢安在繁花嫩叶当中，高兴起来了，呼朋引伴地卖弄清脆的喉咙，唱出宛转的曲子，与轻风流水应和着。牛背上牧童的短笛，这时候也成天嘹亮地响着。\n\n雨是最寻常的，一下就是三两天。可别恼。看，像牛毛，像花针，像细丝，密密地斜织着，人家屋顶上全笼着一层薄烟。树叶儿却绿得发亮，小草儿也青得逼你的眼。傍晚时候，上灯了，一点点黄晕的光，烘托出一片安静而和平的夜。\n\n在乡下，小路上，石桥边，有撑起伞慢慢走着的人，地里还有工作的农民，披着蓑戴着笠。他们的房屋，稀稀疏疏的在雨里静默着。\n\n天上风筝渐渐多了，地上孩子也多了。城里乡下，家家户户，老老小小，也赶趟儿似的，一个个都出来了。舒活舒活筋骨，抖擞抖擞精神，各做各的一份事去。“一年之计在于春”，刚起头儿，有的是工夫，有的是希望。\n\n春天像刚落地的娃娃，从头到脚都是新的，它生长着。春天像小姑娘，花枝招展的，笑着，走着。春天像健壮的青年，有铁一般的胳膊和腰脚，领着我们上前去。', 1, 1, '散文|朱自清', 1, 0, '/static/upload/20201104/d929f4f9fc1e7c7ea65327589be2e99f.jpg', 1604491635, 1604570098, 1604570098);
INSERT INTO `blog_article` VALUES (6, 'MuAdmin', 'MuAdmin真不错', 1, '### MuxAdmin 真不错\n\n真的', 1, 1, 'MuxAdmin|后端模板', 7, 1, '/static/upload/20201105/cab0b764ffd889a34b4d5b49c22fd2ba.jpg', 1604571183, 1604571183, NULL);
COMMIT;

-- ----------------------------
-- Table structure for blog_category
-- ----------------------------
DROP TABLE IF EXISTS `blog_category`;
CREATE TABLE `blog_category` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `state` int NOT NULL DEFAULT '1' COMMENT '1显示0隐藏',
  `sort` int NOT NULL,
  `create_time` int DEFAULT NULL,
  `update_time` int DEFAULT NULL,
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of blog_category
-- ----------------------------
BEGIN;
INSERT INTO `blog_category` VALUES (1, 'PHP', 1, 1, 1604148043, 1604153194, NULL);
INSERT INTO `blog_category` VALUES (2, 'Java', 1, 3, 1604139786, 1604153504, NULL);
INSERT INTO `blog_category` VALUES (3, 'Linux', 1, 6, 1604139870, 1604563173, NULL);
INSERT INTO `blog_category` VALUES (4, 'HTML', 1, 5, 1604153428, 1604554092, NULL);
INSERT INTO `blog_category` VALUES (5, 'Javascript', 1, 4, 1604153442, 1604563039, NULL);
INSERT INTO `blog_category` VALUES (6, 'CSS', 1, 5, 1604153454, 1604568780, 1604568780);
INSERT INTO `blog_category` VALUES (7, 'muxadmin', 1, 1, 1604569034, 1604569034, NULL);
COMMIT;

-- ----------------------------
-- Table structure for blog_web
-- ----------------------------
DROP TABLE IF EXISTS `blog_web`;
CREATE TABLE `blog_web` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `title` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `keyword` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `desc` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `copyright` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `beian` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `state` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `message` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of blog_web
-- ----------------------------
BEGIN;
INSERT INTO `blog_web` VALUES (1, 'PureBlog', '轻量级博客系统', 'PureBlog,blog', '轻量级博客系统', 'Copyright © 2020.Company name All rights reserved.', '赣ICP备00010000号', '0', '服务器正在打瞌睡!');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
