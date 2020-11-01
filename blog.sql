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

 Date: 01/11/2020 14:24:22
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
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of blog_admin
-- ----------------------------
BEGIN;
INSERT INTO `blog_admin` VALUES (1, 'admin', '21232f297a57a5a743894a0e4a801fc3', 1, 1, 1604037851, 1604129747, NULL);
INSERT INTO `blog_admin` VALUES (2, 'test1', 'e10adc3949ba59abbe56e057f20f883e', 0, 1, 1604117061, 1604201076, NULL);
INSERT INTO `blog_admin` VALUES (3, 'test2', '8ad8757baa8564dc136c1e07507f4a98', 0, 1, 1604117061, 1604141955, NULL);
INSERT INTO `blog_admin` VALUES (4, 'test3', '8ad8757baa8564dc136c1e07507f4a98', 1, 1, 1604117153, 1604128325, NULL);
INSERT INTO `blog_admin` VALUES (5, 'test4', '86985e105f79b95d6bc918fb45ec7727', 1, 1, 1604117202, 1604117202, NULL);
INSERT INTO `blog_admin` VALUES (6, 'test6', '4cfad7076129962ee70c36839a1e3e15', 0, 1, 1604117218, 1604150711, NULL);
INSERT INTO `blog_admin` VALUES (7, 'test', '8ad8757baa8564dc136c1e07507f4a98', 0, 1, 1604129783, 1604129783, NULL);
INSERT INTO `blog_admin` VALUES (8, 'test34', '8ad8757baa8564dc136c1e07507f4a98', 0, 1, 1604129792, 1604129792, NULL);
INSERT INTO `blog_admin` VALUES (9, 'test5', '8ad8757baa8564dc136c1e07507f4a98', 0, 1, 1604129801, 1604129801, NULL);
INSERT INTO `blog_admin` VALUES (10, 'abcd', 'e2fc714c4727ee9395f324cd2e7f331f', 0, 1, 1604145115, 1604150699, NULL);
COMMIT;

-- ----------------------------
-- Table structure for blog_article
-- ----------------------------
DROP TABLE IF EXISTS `blog_article`;
CREATE TABLE `blog_article` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `top` int NOT NULL DEFAULT '0',
  `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `state` int NOT NULL DEFAULT '1',
  `u_id` int NOT NULL,
  `tag` varchar(255) NOT NULL,
  `category_id` int NOT NULL,
  `read` int DEFAULT NULL,
  `pic` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `create_time` int DEFAULT NULL,
  `update_time` int DEFAULT NULL,
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of blog_article
-- ----------------------------
BEGIN;
INSERT INTO `blog_article` VALUES (1, 'Hello World', 1, '# PureBlog TODO\n\n\n\n#### 基本模板\n\n- [x] 博客模板\n- [x] 后台模板\n\n---\n\n#### 登录模块\n\n- [x] 后台登录\n- [x] 后台主页\n- [x] 后台退出登录\n\n---\n\n#### 管理员模块\n\n- [x] 管理员列表\n- [x] 管理员权限修饰\n- [x] 管理员添加\n- [x] 管理员修改\n- [x] 管理员分页\n- [x] 管理员搜索\n- [x] 管理员删除\n- [x] 管理员状态修改\n\n---\n\n#### 栏目模块\n\n- [x] 栏目列表\n- [x] 栏目添加\n- [x] 栏目修改\n- [x] 栏目删除\n- [x] 栏目分页\n- [x] 栏目排序\n- [x] 栏目状态修改\n\n---\n\n#### 文章模块\n\n- [x] 文章列表\n- [ ] 文章添加\n- [ ] 文章修改\n- [ ] 文章删除\n- [ ] 文章分页\n- [ ] TODO\n\n', 1, 1, 'js|a', 1, 99, '/static/upload/test.jpg', 1, NULL, NULL);
INSERT INTO `blog_article` VALUES (2, 'test', 1, '# Test', 1, 1, 'test|html', 1, NULL, '/static/upload/20201101/28a8d5c18b7583ffc531cf01efa22a2b.jpg', 1604192552, NULL, NULL);
INSERT INTO `blog_article` VALUES (3, 'Test', 1, '##  如何优雅的写php\n\n```php\n<?php\necho \"你妈死了\"\n\n```', 1, 1, 'Test', 3, NULL, '/static/upload/20201101/a73f65f57ee9050f9a00e1d67b3b08dc.png', 1604211524, 1604211540, NULL);
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
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of blog_category
-- ----------------------------
BEGIN;
INSERT INTO `blog_category` VALUES (1, 'PHP', 1, 1, 1604148043, 1604153194, NULL);
INSERT INTO `blog_category` VALUES (2, 'Java', 1, 3, 1604139786, 1604153504, NULL);
INSERT INTO `blog_category` VALUES (3, 'Linux', 1, 2, 1604139870, 1604153485, NULL);
INSERT INTO `blog_category` VALUES (4, 'Test1', 1, 5, 1604153428, 1604153518, NULL);
INSERT INTO `blog_category` VALUES (5, 'Test2', 1, 3, 1604153442, 1604153493, NULL);
INSERT INTO `blog_category` VALUES (6, 'Test3', 1, 5, 1604153454, 1604153454, NULL);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
