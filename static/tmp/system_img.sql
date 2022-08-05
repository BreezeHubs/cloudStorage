/*
 Navicat Premium Data Transfer

 Source Server         : i026
 Source Server Type    : MySQL
 Source Server Version : 50734
 Source Host           : 119.23.201.145:3306
 Source Schema         : ipxe_local

 Target Server Type    : MySQL
 Target Server Version : 50734
 File Encoding         : 65001

 Date: 16/12/2021 14:59:38
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for system_img
-- ----------------------------
DROP TABLE IF EXISTS `system_img`;
CREATE TABLE `system_img`  (
  `id` int(11) NOT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '系统镜像版本',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of system_img
-- ----------------------------
INSERT INTO `system_img` VALUES (1, 'nfs');
INSERT INTO `system_img` VALUES (2, '独立镜像-N卡');

SET FOREIGN_KEY_CHECKS = 1;
