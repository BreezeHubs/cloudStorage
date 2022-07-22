/*
 Navicat Premium Data Transfer

 Source Server         : mysql-master
 Source Server Type    : MySQL
 Source Server Version : 80029
 Source Host           : localhost:3306
 Source Schema         : tbl_file

 Target Server Type    : MySQL
 Target Server Version : 80029
 File Encoding         : 65001

 Date: 22/07/2022 18:08:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for tbl_file
-- ----------------------------
DROP TABLE IF EXISTS `tbl_file`;
CREATE TABLE `tbl_file`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `file_hash` char(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint NOT NULL DEFAULT 0 COMMENT '文件大小',
  `file_addr` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `create_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` int NOT NULL DEFAULT 0 COMMENT '状态(可用/禁用/已删除等)',
  `ext1` int NULL DEFAULT 0 COMMENT '扩展字段1',
  `ext2` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '扩展字段2',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE,
  UNIQUE INDEX `idx_hash`(`file_hash`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '文件表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of tbl_file
-- ----------------------------
INSERT INTO `tbl_file` VALUES (1, '93e46ae4c7d9087e7070dde35c79197eae52b7e8', '1.png', 33288, './static/tmp/1.png', '2022-07-19 07:45:57', '2022-07-21 09:09:36', 1, 0, NULL);
INSERT INTO `tbl_file` VALUES (2, '1ee37dc2e1790ccd814c955f4e080bbd9d134a74', 'vulkan-1.dll', 696664, './static/tmp/vulkan-1.dll', '2022-07-21 06:51:01', '2022-07-21 06:51:01', 1, 0, NULL);
INSERT INTO `tbl_file` VALUES (3, 'acddd220fd3fd9d869652b71a8c685c3c5289f86', '1.3.1驱动程序设置.jpg', 395643, './static/tmp/1.3.1驱动程序设置.jpg', '2022-07-21 06:51:15', '2022-07-21 06:51:15', 1, 0, NULL);
INSERT INTO `tbl_file` VALUES (4, '9934967c4cb31f0313142de893918a9636fee175', '本地调度系统-设计说明书.docx', 23249, './static/tmp/本地调度系统-设计说明书.docx', '2022-07-21 09:15:35', '2022-07-21 09:15:35', 1, 0, NULL);
INSERT INTO `tbl_file` VALUES (5, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 28098, './static/file/771188497ae76d762c8ebdbf4bb13804dd25411d.docx', '2022-07-22 06:12:14', '2022-07-22 06:12:14', 1, 0, NULL);
INSERT INTO `tbl_file` VALUES (9, 'dfa39cac093a7a9c94d25130671ec474d51a2995', '代码资料.zip', 132489256, '', '2022-07-22 09:30:51', '2022-07-22 09:30:51', 1, 0, NULL);

-- ----------------------------
-- Table structure for tbl_user
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user`;
CREATE TABLE `tbl_user`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `user_pwd` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户密码',
  `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '邮箱',
  `phone` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '手机号',
  `email_vaildated` tinyint(1) NULL DEFAULT 0 COMMENT '邮箱是否验证',
  `phone_vaildated` tinyint(1) NULL DEFAULT 0 COMMENT '手机号是否验证',
  `create_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  `update_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间',
  `profile` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '用户属性',
  `status` int NOT NULL DEFAULT 0 COMMENT '状态(可用/禁用/已删除等)',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of tbl_user
-- ----------------------------
INSERT INTO `tbl_user` VALUES (5, 'admin', '1a48f352ffe1e15d1c4611b22615df0092ec194c', '', '', 0, 0, '2022-07-20 08:51:28', '2022-07-20 08:51:28', NULL, 0);
INSERT INTO `tbl_user` VALUES (7, 'user01', 'a4edd8abbab15d6ebebc8b3d36443a0803d509c2', '', '', 0, 0, '2022-07-22 02:16:26', '2022-07-22 02:16:26', NULL, 0);

-- ----------------------------
-- Table structure for tbl_user_file
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user_file`;
CREATE TABLE `tbl_user_file`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户id',
  `file_hash` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint NULL DEFAULT 0 COMMENT '文件大小',
  `create_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` int NOT NULL DEFAULT 0 COMMENT '状态(可用/禁用/已删除等)',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户文件表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of tbl_user_file
-- ----------------------------
INSERT INTO `tbl_user_file` VALUES (1, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 28098, '2022-07-22 06:30:08', '2022-07-22 14:30:05', 1);
INSERT INTO `tbl_user_file` VALUES (2, 7, 'acddd220fd3fd9d869652b71a8c685c3c5289f86', 'shezhishezhi.jpg', 395643, '2022-07-22 07:00:01', '2022-07-22 07:00:01', 1);
INSERT INTO `tbl_user_file` VALUES (3, 7, 'dfa39cac093a7a9c94d25130671ec474d51a2995', '代码资料.zip', 132489256, '2022-07-22 09:30:51', '2022-07-22 09:30:51', 1);
INSERT INTO `tbl_user_file` VALUES (4, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 132489256, '2022-07-22 09:33:32', '2022-07-22 09:33:32', 1);
INSERT INTO `tbl_user_file` VALUES (5, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 132489256, '2022-07-22 09:42:11', '2022-07-22 09:42:11', 1);
INSERT INTO `tbl_user_file` VALUES (6, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 132489256, '2022-07-22 09:47:00', '2022-07-22 09:47:00', 1);
INSERT INTO `tbl_user_file` VALUES (7, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 132489256, '2022-07-22 09:49:04', '2022-07-22 09:49:04', 1);
INSERT INTO `tbl_user_file` VALUES (8, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 132489256, '2022-07-22 09:49:16', '2022-07-22 09:49:16', 1);
INSERT INTO `tbl_user_file` VALUES (9, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 132489256, '2022-07-22 09:50:11', '2022-07-22 09:50:11', 1);
INSERT INTO `tbl_user_file` VALUES (10, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 132489256, '2022-07-22 09:50:34', '2022-07-22 09:50:34', 1);
INSERT INTO `tbl_user_file` VALUES (11, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 132489256, '2022-07-22 09:52:20', '2022-07-22 09:52:20', 1);
INSERT INTO `tbl_user_file` VALUES (12, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 132489256, '2022-07-22 09:53:11', '2022-07-22 09:53:11', 1);
INSERT INTO `tbl_user_file` VALUES (13, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 132489256, '2022-07-22 09:57:45', '2022-07-22 09:57:45', 1);
INSERT INTO `tbl_user_file` VALUES (14, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 28098, '2022-07-22 09:59:32', '2022-07-22 09:59:32', 1);
INSERT INTO `tbl_user_file` VALUES (15, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 28098, '2022-07-22 10:01:48', '2022-07-22 10:01:48', 1);
INSERT INTO `tbl_user_file` VALUES (16, 7, '771188497ae76d762c8ebdbf4bb13804dd25411d', 'snmp使用分析.docx', 28098, '2022-07-22 10:03:04', '2022-07-22 10:03:04', 1);

-- ----------------------------
-- Table structure for tbl_user_token
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user_token`;
CREATE TABLE `tbl_user_token`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL DEFAULT '' COMMENT '用户id',
  `user_token` char(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'token',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 27 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户token表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of tbl_user_token
-- ----------------------------
INSERT INTO `tbl_user_token` VALUES (27, 7, '26ba9ed1aae51dfeb32538f4697f7e6b62da3eaf');

SET FOREIGN_KEY_CHECKS = 1;
