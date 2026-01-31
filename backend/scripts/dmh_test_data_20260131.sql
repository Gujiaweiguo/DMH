mysqldump: [Warning] Using a password on the command line interface can be insecure.
-- MySQL dump 10.13  Distrib 8.0.44, for Linux (x86_64)
--
-- Host: localhost    Database: dmh
-- ------------------------------------------------------
-- Server version	8.0.44

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `audit_logs`
--

DROP TABLE IF EXISTS `audit_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `audit_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint DEFAULT NULL COMMENT '操作用户ID',
  `username` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '操作用户名',
  `action` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作类型',
  `resource` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '操作资源',
  `resource_id` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '资源ID',
  `details` text COLLATE utf8mb4_unicode_ci COMMENT '操作详情',
  `client_ip` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '客户端IP',
  `user_agent` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '用户代理',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'success' COMMENT '操作状态: success/failed',
  `error_msg` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '错误信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_action` (`action`),
  KEY `idx_resource` (`resource`),
  KEY `idx_resource_id` (`resource_id`),
  KEY `idx_client_ip` (`client_ip`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `audit_logs_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作审计日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `audit_logs`
--

LOCK TABLES `audit_logs` WRITE;
/*!40000 ALTER TABLE `audit_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `audit_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `brand_assets`
--

DROP TABLE IF EXISTS `brand_assets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `brand_assets` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '素材ID',
  `brand_id` bigint NOT NULL COMMENT '品牌ID',
  `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '素材名称',
  `type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '素材类型: image/video/document',
  `category` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '素材分类',
  `tags` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '素材标签',
  `file_url` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件URL',
  `file_size` bigint DEFAULT '0' COMMENT '文件大小（字节）',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '素材描述',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '状态: active/disabled',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_brand_id` (`brand_id`),
  KEY `idx_type` (`type`),
  KEY `idx_category` (`category`),
  KEY `idx_status` (`status`),
  CONSTRAINT `brand_assets_ibfk_1` FOREIGN KEY (`brand_id`) REFERENCES `brands` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌素材表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `brand_assets`
--

LOCK TABLES `brand_assets` WRITE;
/*!40000 ALTER TABLE `brand_assets` DISABLE KEYS */;
/*!40000 ALTER TABLE `brand_assets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `brands`
--

DROP TABLE IF EXISTS `brands`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `brands` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '品牌ID',
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '品牌名称',
  `logo` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '品牌Logo',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '品牌描述',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '品牌状态: active/disabled',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `brands`
--

LOCK TABLES `brands` WRITE;
/*!40000 ALTER TABLE `brands` DISABLE KEYS */;
INSERT INTO `brands` VALUES (1,'星巴克咖啡','https://via.placeholder.com/150','全球知名咖啡连锁品牌，提供优质咖啡和温馨环境','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(2,'麦当劳','https://via.placeholder.com/150','全球最大的快餐连锁企业','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(3,'华为科技','https://via.placeholder.com/150','全球领先的ICT基础设施和智能终端提供商','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(4,'小米科技','https://via.placeholder.com/150','专注于智能硬件和电子产品开发的互联网公司','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(5,'耐克运动','https://via.placeholder.com/150','全球著名的运动品牌','active','2026-01-31 10:48:11','2026-01-31 10:48:11');
/*!40000 ALTER TABLE `brands` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `campaigns`
--

DROP TABLE IF EXISTS `campaigns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `campaigns` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '活动ID',
  `brand_id` bigint DEFAULT '0' COMMENT '品牌ID',
  `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '活动名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '活动描述',
  `form_fields` json DEFAULT NULL COMMENT '动态表单字段配置',
  `reward_rule` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '奖励规则（比例或固定金额）',
  `start_time` datetime NOT NULL COMMENT '活动开始时间',
  `end_time` datetime NOT NULL COMMENT '活动结束时间',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '活动状态: active/paused/ended',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_end_time` (`end_time`),
  KEY `idx_brand_id` (`brand_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='营销活动表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `campaigns`
--

LOCK TABLES `campaigns` WRITE;
/*!40000 ALTER TABLE `campaigns` DISABLE KEYS */;
INSERT INTO `campaigns` VALUES (1,1,'星巴克新年促销活动','2026新年大促，分享推广有礼','\"姓名,手机号,地址\"',15.00,'2026-01-01 00:00:00','2026-12-31 23:59:59','active','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL),(2,2,'麦当劳超值套餐推广','推广麦当劳超值套餐，赚取佣金','\"姓名,手机号,备注\"',12.00,'2026-01-15 00:00:00','2026-03-31 23:59:59','active','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL),(3,3,'华为新品体验活动','体验华为最新产品，分享获得奖励','\"姓名,手机号,收货地址\"',20.00,'2026-01-20 00:00:00','2026-06-30 23:59:59','active','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL),(4,1,'星巴克春季促销','春季特惠活动暂时调整','\"姓名,手机号\"',10.00,'2026-03-01 00:00:00','2026-05-31 23:59:59','paused','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL),(5,2,'双十一促销活动','2025双十一已结束','\"姓名,手机号,地址\"',18.00,'2025-11-01 00:00:00','2025-11-11 23:59:59','ended','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL);
/*!40000 ALTER TABLE `campaigns` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `external_orders`
--

DROP TABLE IF EXISTS `external_orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `external_orders` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `order_id` bigint NOT NULL COMMENT 'DMH订单ID',
  `campaign_id` bigint NOT NULL COMMENT '活动ID',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
  `form_data` text COLLATE utf8mb4_unicode_ci COMMENT '表单数据JSON',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '订单金额',
  `pay_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '支付状态',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `synced_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_id` (`order_id`),
  KEY `idx_phone` (`phone`),
  KEY `idx_pay_status` (`pay_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='外部订单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `external_orders`
--

LOCK TABLES `external_orders` WRITE;
/*!40000 ALTER TABLE `external_orders` DISABLE KEYS */;
/*!40000 ALTER TABLE `external_orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `external_rewards`
--

DROP TABLE IF EXISTS `external_rewards`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `external_rewards` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `reward_id` bigint NOT NULL COMMENT 'DMH奖励ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `order_id` bigint NOT NULL COMMENT '订单ID',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '奖励金额',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '奖励状态',
  `settled_at` datetime DEFAULT NULL COMMENT '结算时间',
  `synced_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reward_id` (`reward_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='外部奖励表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `external_rewards`
--

LOCK TABLES `external_rewards` WRITE;
/*!40000 ALTER TABLE `external_rewards` DISABLE KEYS */;
/*!40000 ALTER TABLE `external_rewards` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `login_attempts`
--

DROP TABLE IF EXISTS `login_attempts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `login_attempts` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `user_id` bigint DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '尝试登录的用户名',
  `client_ip` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '客户端IP',
  `user_agent` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '用户代理',
  `success` tinyint(1) NOT NULL COMMENT '是否成功',
  `fail_reason` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '失败原因',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_client_ip` (`client_ip`),
  KEY `idx_success` (`success`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `login_attempts_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录尝试记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `login_attempts`
--

LOCK TABLES `login_attempts` WRITE;
/*!40000 ALTER TABLE `login_attempts` DISABLE KEYS */;
/*!40000 ALTER TABLE `login_attempts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `menus`
--

DROP TABLE IF EXISTS `menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `menus` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单名称',
  `code` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单编码',
  `path` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '菜单路径',
  `icon` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '菜单图标',
  `parent_id` bigint DEFAULT NULL COMMENT '父菜单ID',
  `sort` int DEFAULT '0' COMMENT '排序',
  `type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'menu' COMMENT '类型: menu/button',
  `platform` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'admin' COMMENT '平台: admin/h5',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '状态: active/disabled',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_platform` (`platform`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `menus`
--

LOCK TABLES `menus` WRITE;
/*!40000 ALTER TABLE `menus` DISABLE KEYS */;
INSERT INTO `menus` VALUES (1,'首页','dashboard','/dashboard','home',0,1,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(2,'用户管理','users','/users','users',0,2,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(3,'品牌管理','brands','/brands','shopping-bag',0,3,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(4,'活动管理','campaigns','/campaigns','gift',0,4,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(5,'订单管理','orders','/orders','file-text',0,5,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(6,'分销管理','distributors','/distributors','share-2',0,6,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(7,'数据统计','analytics','/analytics','bar-chart',0,7,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(8,'系统设置','settings','/settings','settings',0,8,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(9,'用户列表','users.list','/users/list','user',2,1,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(10,'品牌列表','brands.list','/brands/list','list',3,1,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(11,'活动列表','campaigns.list','/campaigns/list','list',4,1,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(12,'活动编辑','campaigns.edit','/campaigns/edit','edit',4,2,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(13,'订单列表','orders.list','/orders/list','list',5,1,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(14,'分销商列表','distributors.list','/distributors/list','list',6,1,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(15,'推广奖励','distributors.rewards','/distributors/rewards','award',6,2,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(16,'活动统计','analytics.campaigns','/analytics/campaigns','pie-chart',7,1,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(17,'用户统计','analytics.users','/analytics/users','users',7,2,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(18,'基本设置','settings.basic','/settings/basic','settings',8,1,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11'),(19,'权限设置','settings.permissions','/settings/permissions','lock',8,2,'menu','admin','active','2026-01-31 10:48:11','2026-01-31 10:48:11');
/*!40000 ALTER TABLE `menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `campaign_id` bigint NOT NULL COMMENT '活动ID',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户手机号',
  `form_data` json DEFAULT NULL COMMENT '表单数据',
  `referrer_id` bigint DEFAULT '0' COMMENT '推荐人ID',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '订单状态: pending/paid/cancelled',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '订单金额',
  `pay_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'unpaid' COMMENT '支付状态: unpaid/paid/refunded',
  `trade_no` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '交易流水号',
  `sync_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '同步状态: pending/synced/failed',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_campaign_phone` (`campaign_id`,`phone`,`deleted_at`) COMMENT '同一活动同一手机号只能报名一次',
  KEY `idx_campaign_id` (`campaign_id`),
  KEY `idx_phone` (`phone`),
  KEY `idx_referrer_id` (`referrer_id`),
  KEY `idx_status` (`status`),
  KEY `idx_pay_status` (`pay_status`),
  KEY `idx_sync_status` (`sync_status`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
INSERT INTO `orders` VALUES (1,1,'13900001001','{\"name\": \"张三\", \"phone\": \"13900001001\", \"address\": \"北京市朝阳区\"}',7,'paid',99.00,'paid','TX20260131001','synced','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL),(2,1,'13900001002','{\"name\": \"李四\", \"phone\": \"13900001002\", \"address\": \"上海市浦东新区\"}',8,'paid',199.00,'paid','TX20260131002','synced','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL),(3,2,'13900001003','{\"name\": \"王五\", \"note\": \"无糖\", \"phone\": \"13900001003\"}',7,'paid',59.00,'paid','TX20260131003','synced','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL),(4,3,'13900001004','{\"name\": \"赵六\", \"phone\": \"13900001004\", \"address\": \"广州市天河区\"}',7,'paid',299.00,'paid','TX20260131004','synced','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL),(5,1,'13900001005','{\"name\": \"孙七\", \"phone\": \"13900001005\", \"address\": \"深圳市南山区\"}',0,'pending',99.00,'unpaid','','pending','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL),(6,2,'13900001006','{\"name\": \"周八\", \"note\": \"加芝士\", \"phone\": \"13900001006\"}',0,'pending',59.00,'unpaid','','pending','2026-01-31 10:48:11','2026-01-31 10:48:11',NULL);
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `password_histories`
--

DROP TABLE IF EXISTS `password_histories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `password_histories` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '历史ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `password_histories_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密码历史记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `password_histories`
--

LOCK TABLES `password_histories` WRITE;
/*!40000 ALTER TABLE `password_histories` DISABLE KEYS */;
/*!40000 ALTER TABLE `password_histories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `password_policies`
--

DROP TABLE IF EXISTS `password_policies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `password_policies` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '策略ID',
  `min_length` int NOT NULL DEFAULT '8' COMMENT '最小长度',
  `require_uppercase` tinyint(1) NOT NULL DEFAULT '1' COMMENT '需要大写字母',
  `require_lowercase` tinyint(1) NOT NULL DEFAULT '1' COMMENT '需要小写字母',
  `require_numbers` tinyint(1) NOT NULL DEFAULT '1' COMMENT '需要数字',
  `require_special_chars` tinyint(1) NOT NULL DEFAULT '1' COMMENT '需要特殊字符',
  `max_age` int NOT NULL DEFAULT '90' COMMENT '密码最大有效期（天）',
  `history_count` int NOT NULL DEFAULT '5' COMMENT '历史密码记录数量',
  `max_login_attempts` int NOT NULL DEFAULT '5' COMMENT '最大登录尝试次数',
  `lockout_duration` int NOT NULL DEFAULT '30' COMMENT '锁定时长（分钟）',
  `session_timeout` int NOT NULL DEFAULT '480' COMMENT '会话超时时间（分钟）',
  `max_concurrent_sessions` int NOT NULL DEFAULT '3' COMMENT '最大并发会话数',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密码策略配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `password_policies`
--

LOCK TABLES `password_policies` WRITE;
/*!40000 ALTER TABLE `password_policies` DISABLE KEYS */;
INSERT INTO `password_policies` VALUES (1,8,1,1,1,1,90,5,5,30,480,3,'2026-01-31 10:40:22','2026-01-31 10:40:22');
/*!40000 ALTER TABLE `password_policies` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permissions`
--

DROP TABLE IF EXISTS `permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `permissions` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '权限ID',
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限名称',
  `code` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限编码',
  `resource` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型',
  `action` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作类型',
  `description` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '权限描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_resource` (`resource`),
  KEY `idx_action` (`action`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permissions`
--

LOCK TABLES `permissions` WRITE;
/*!40000 ALTER TABLE `permissions` DISABLE KEYS */;
INSERT INTO `permissions` VALUES (1,'查看用户列表','user:list','user','list','查看所有用户信息','2026-01-31 10:48:11'),(2,'创建用户','user:create','user','create','创建新用户','2026-01-31 10:48:11'),(3,'更新用户','user:update','user','update','更新用户信息','2026-01-31 10:48:11'),(4,'删除用户','user:delete','user','delete','删除用户','2026-01-31 10:48:11'),(5,'查看品牌列表','brand:list','brand','list','查看所有品牌','2026-01-31 10:48:11'),(6,'创建品牌','brand:create','brand','create','创建新品牌','2026-01-31 10:48:11'),(7,'更新品牌','brand:update','brand','update','更新品牌信息','2026-01-31 10:48:11'),(8,'查看活动列表','campaign:list','campaign','list','查看所有营销活动','2026-01-31 10:48:11'),(9,'创建活动','campaign:create','campaign','create','创建新活动','2026-01-31 10:48:11'),(10,'更新活动','campaign:update','campaign','update','更新活动信息','2026-01-31 10:48:11'),(11,'删除活动','campaign:delete','campaign','delete','删除活动','2026-01-31 10:48:11'),(12,'查看订单列表','order:list','order','list','查看所有订单','2026-01-31 10:48:11'),(13,'导出订单','order:export','order','export','导出订单数据','2026-01-31 10:48:11');
/*!40000 ALTER TABLE `permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rewards`
--

DROP TABLE IF EXISTS `rewards`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rewards` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '奖励ID',
  `user_id` bigint NOT NULL COMMENT '用户ID（推荐人）',
  `order_id` bigint NOT NULL COMMENT '关联订单ID',
  `campaign_id` bigint NOT NULL COMMENT '活动ID',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '奖励金额',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '奖励状态: pending/settled/cancelled',
  `settled_at` datetime DEFAULT NULL COMMENT '结算时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_campaign_id` (`campaign_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='奖励记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rewards`
--

LOCK TABLES `rewards` WRITE;
/*!40000 ALTER TABLE `rewards` DISABLE KEYS */;
INSERT INTO `rewards` VALUES (1,7,1,1,15.00,'settled','2026-01-30 10:00:00','2026-01-30 09:30:00','2026-01-30 10:00:00'),(2,7,2,1,15.00,'settled','2026-01-30 11:00:00','2026-01-30 10:30:00','2026-01-30 11:00:00'),(3,8,3,2,12.00,'settled','2026-01-30 12:00:00','2026-01-30 11:30:00','2026-01-30 12:00:00'),(4,7,4,3,20.00,'settled','2026-01-30 13:00:00','2026-01-30 12:30:00','2026-01-30 13:00:00'),(5,7,5,1,15.00,'pending',NULL,'2026-01-31 08:00:00','2026-01-31 08:00:00'),(6,8,6,2,12.00,'pending',NULL,'2026-01-31 09:00:00','2026-01-31 09:00:00');
/*!40000 ALTER TABLE `rewards` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_menus`
--

DROP TABLE IF EXISTS `role_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_menus` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `menu_id` bigint NOT NULL COMMENT '菜单ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menu` (`role_id`,`menu_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_menus`
--

LOCK TABLES `role_menus` WRITE;
/*!40000 ALTER TABLE `role_menus` DISABLE KEYS */;
INSERT INTO `role_menus` VALUES (1,1,1,'2026-01-31 10:48:11'),(2,1,2,'2026-01-31 10:48:11'),(3,1,3,'2026-01-31 10:48:11'),(4,1,4,'2026-01-31 10:48:11'),(5,1,5,'2026-01-31 10:48:11'),(6,1,6,'2026-01-31 10:48:11'),(7,1,7,'2026-01-31 10:48:11'),(8,1,8,'2026-01-31 10:48:11'),(9,1,9,'2026-01-31 10:48:11'),(10,1,10,'2026-01-31 10:48:11'),(11,1,11,'2026-01-31 10:48:11'),(12,1,12,'2026-01-31 10:48:11'),(13,1,13,'2026-01-31 10:48:11'),(14,1,14,'2026-01-31 10:48:11'),(15,1,15,'2026-01-31 10:48:11'),(16,1,16,'2026-01-31 10:48:11'),(17,1,17,'2026-01-31 10:48:11'),(18,1,18,'2026-01-31 10:48:11'),(19,1,19,'2026-01-31 10:48:11'),(20,2,1,'2026-01-31 10:48:11'),(21,2,3,'2026-01-31 10:48:11'),(22,2,4,'2026-01-31 10:48:11'),(23,2,5,'2026-01-31 10:48:11'),(24,2,6,'2026-01-31 10:48:11'),(25,2,7,'2026-01-31 10:48:11'),(26,2,10,'2026-01-31 10:48:11'),(27,2,11,'2026-01-31 10:48:11'),(28,2,12,'2026-01-31 10:48:11'),(29,2,13,'2026-01-31 10:48:11'),(30,2,14,'2026-01-31 10:48:11'),(31,2,15,'2026-01-31 10:48:11'),(32,2,16,'2026-01-31 10:48:11');
/*!40000 ALTER TABLE `role_menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_permissions`
--

DROP TABLE IF EXISTS `role_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_permissions` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `permission_id` bigint NOT NULL COMMENT '权限ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`,`permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_permissions`
--

LOCK TABLES `role_permissions` WRITE;
/*!40000 ALTER TABLE `role_permissions` DISABLE KEYS */;
INSERT INTO `role_permissions` VALUES (1,1,1,'2026-01-31 10:48:11'),(2,1,2,'2026-01-31 10:48:11'),(3,1,3,'2026-01-31 10:48:11'),(4,1,4,'2026-01-31 10:48:11'),(5,1,5,'2026-01-31 10:48:11'),(6,1,6,'2026-01-31 10:48:11'),(7,1,7,'2026-01-31 10:48:11'),(8,1,8,'2026-01-31 10:48:11'),(9,1,9,'2026-01-31 10:48:11'),(10,1,10,'2026-01-31 10:48:11'),(11,1,11,'2026-01-31 10:48:11'),(12,1,12,'2026-01-31 10:48:11'),(13,1,13,'2026-01-31 10:48:11'),(14,2,5,'2026-01-31 10:48:11'),(15,2,6,'2026-01-31 10:48:11'),(16,2,7,'2026-01-31 10:48:11'),(17,2,8,'2026-01-31 10:48:11'),(18,2,9,'2026-01-31 10:48:11'),(19,2,10,'2026-01-31 10:48:11'),(20,2,11,'2026-01-31 10:48:11'),(21,2,12,'2026-01-31 10:48:11'),(22,2,13,'2026-01-31 10:48:11');
/*!40000 ALTER TABLE `role_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色编码',
  `description` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '角色描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (1,'平台管理员','platform_admin','拥有系统所有权限的超级管理员','2026-01-31 10:40:22','2026-01-31 10:40:22'),(2,'参与者','participant','普通用户，可参与活动','2026-01-31 10:40:22','2026-01-31 10:40:22'),(3,'匿名用户','anonymous','未登录的访客用户','2026-01-31 10:40:22','2026-01-31 10:40:22');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `security_events`
--

DROP TABLE IF EXISTS `security_events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `security_events` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '事件ID',
  `event_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '事件类型',
  `severity` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '严重程度: low/medium/high/critical',
  `user_id` bigint DEFAULT NULL COMMENT '相关用户ID',
  `username` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '相关用户名',
  `client_ip` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '客户端IP',
  `user_agent` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '用户代理',
  `description` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '事件描述',
  `details` json DEFAULT NULL COMMENT '事件详情',
  `handled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已处理',
  `handled_by` bigint DEFAULT NULL COMMENT '处理人ID',
  `handled_at` datetime DEFAULT NULL COMMENT '处理时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_event_type` (`event_type`),
  KEY `idx_severity` (`severity`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_client_ip` (`client_ip`),
  KEY `idx_handled` (`handled`),
  KEY `idx_created_at` (`created_at`),
  KEY `handled_by` (`handled_by`),
  CONSTRAINT `security_events_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  CONSTRAINT `security_events_ibfk_2` FOREIGN KEY (`handled_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='安全事件记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `security_events`
--

LOCK TABLES `security_events` WRITE;
/*!40000 ALTER TABLE `security_events` DISABLE KEYS */;
/*!40000 ALTER TABLE `security_events` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sync_logs`
--

DROP TABLE IF EXISTS `sync_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sync_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `order_id` bigint NOT NULL COMMENT '订单ID',
  `sync_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '同步类型: order/reward',
  `sync_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '同步状态: pending/synced/failed',
  `attempts` int NOT NULL DEFAULT '0' COMMENT '尝试次数',
  `error_msg` text COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `synced_at` datetime DEFAULT NULL COMMENT '同步成功时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_sync_status` (`sync_status`),
  KEY `idx_sync_type` (`sync_type`),
  KEY `idx_synced_at` (`synced_at`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='同步日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sync_logs`
--

LOCK TABLES `sync_logs` WRITE;
/*!40000 ALTER TABLE `sync_logs` DISABLE KEYS */;
INSERT INTO `sync_logs` VALUES (1,1,'order','synced',1,NULL,'2026-01-30 10:05:00','2026-01-30 10:00:00','2026-01-30 10:05:00'),(2,2,'order','synced',1,NULL,'2026-01-30 11:05:00','2026-01-30 11:00:00','2026-01-30 11:05:00'),(3,5,'order','failed',3,'外部系统接口超时',NULL,'2026-01-31 08:00:00','2026-01-31 08:05:00'),(4,1,'reward','synced',1,NULL,'2026-01-30 10:10:00','2026-01-30 10:00:00','2026-01-30 10:10:00'),(5,5,'reward','pending',0,NULL,NULL,'2026-01-31 08:00:00','2026-01-31 08:00:00');
/*!40000 ALTER TABLE `sync_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_balances`
--

DROP TABLE IF EXISTS `user_balances`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_balances` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '余额ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `balance` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '当前余额',
  `total_reward` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '累计奖励',
  `version` bigint NOT NULL DEFAULT '0' COMMENT '版本号（乐观锁）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户余额表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_balances`
--

LOCK TABLES `user_balances` WRITE;
/*!40000 ALTER TABLE `user_balances` DISABLE KEYS */;
INSERT INTO `user_balances` VALUES (1,1,0.00,0.00,0,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(2,2,1250.50,3250.50,0,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(3,3,890.00,1890.00,0,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(4,4,0.00,0.00,0,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(5,5,0.00,0.00,0,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(6,6,0.00,0.00,0,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(7,7,680.30,2480.30,0,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(8,8,320.75,1200.75,0,'2026-01-31 10:48:11','2026-01-31 10:48:11');
/*!40000 ALTER TABLE `user_balances` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_brands`
--

DROP TABLE IF EXISTS `user_brands`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_brands` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `brand_id` bigint NOT NULL COMMENT '品牌ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_brand` (`user_id`,`brand_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_brand_id` (`brand_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户品牌关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_brands`
--

LOCK TABLES `user_brands` WRITE;
/*!40000 ALTER TABLE `user_brands` DISABLE KEYS */;
INSERT INTO `user_brands` VALUES (1,2,1,'2026-01-31 10:48:11'),(2,3,2,'2026-01-31 10:48:11'),(3,7,1,'2026-01-31 10:48:11'),(4,7,2,'2026-01-31 10:48:11'),(5,8,1,'2026-01-31 10:48:11');
/*!40000 ALTER TABLE `user_brands` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_roles`
--

DROP TABLE IF EXISTS `user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_roles`
--

LOCK TABLES `user_roles` WRITE;
/*!40000 ALTER TABLE `user_roles` DISABLE KEYS */;
INSERT INTO `user_roles` VALUES (1,1,1,'2026-01-31 10:48:11'),(2,2,2,'2026-01-31 10:48:11'),(3,3,2,'2026-01-31 10:48:11'),(4,4,3,'2026-01-31 10:48:11'),(5,5,3,'2026-01-31 10:48:11'),(6,6,3,'2026-01-31 10:48:11'),(7,7,3,'2026-01-31 10:48:11'),(8,8,3,'2026-01-31 10:48:11');
/*!40000 ALTER TABLE `user_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_sessions`
--

DROP TABLE IF EXISTS `user_sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_sessions` (
  `id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '会话ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `client_ip` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '客户端IP',
  `user_agent` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '用户代理',
  `login_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `last_active_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后活跃时间',
  `expires_at` datetime NOT NULL COMMENT '过期时间',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '状态: active/expired/revoked',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_last_active_at` (`last_active_at`),
  KEY `idx_expires_at` (`expires_at`),
  KEY `idx_status` (`status`),
  CONSTRAINT `user_sessions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户会话记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_sessions`
--

LOCK TABLES `user_sessions` WRITE;
/*!40000 ALTER TABLE `user_sessions` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码（bcrypt加密）',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '邮箱',
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '头像URL',
  `real_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '真实姓名',
  `role` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'participant' COMMENT '用户角色: platform_admin/participant',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '用户状态: active/disabled/locked',
  `login_attempts` int DEFAULT '0' COMMENT '登录尝试次数',
  `locked_until` datetime DEFAULT NULL COMMENT '锁定到期时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_phone` (`phone`),
  KEY `idx_role` (`role`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'admin','$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm','13800000001','admin@dmh.com','','系统管理员','platform_admin','active',0,NULL,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(2,'brand_manager','$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm','13800000002','brand@dmh.com','','品牌经理','brand_admin','active',0,NULL,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(3,'brand_admin','$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm','13800000003','brand_admin@dmh.com','','品牌管理员','brand_admin','active',0,NULL,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(4,'user001','$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm','13800000004','user001@dmh.com','','张三','participant','active',0,NULL,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(5,'user002','$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm','13800000005','user002@dmh.com','','李四','participant','active',0,NULL,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(6,'user003','$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm','13800000006','user003@dmh.com','','王五','participant','active',0,NULL,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(7,'distributor001','$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm','13800000007','dist001@dmh.com','','分销商一','participant','active',0,NULL,'2026-01-31 10:48:11','2026-01-31 10:48:11'),(8,'distributor002','$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm','13800000008','dist002@dmh.com','','分销商二','participant','active',0,NULL,'2026-01-31 10:48:11','2026-01-31 10:48:11');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `withdrawals`
--

DROP TABLE IF EXISTS `withdrawals`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `withdrawals` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '提现ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `amount` decimal(10,2) NOT NULL COMMENT '提现金额',
  `bank_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '银行名称',
  `bank_account` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '银行账号',
  `account_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '账户名称',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '提现状态: pending/approved/rejected',
  `remark` text COLLATE utf8mb4_unicode_ci COMMENT '备注/拒绝原因',
  `approved_by` bigint DEFAULT NULL COMMENT '审核人ID',
  `approved_at` datetime DEFAULT NULL COMMENT '审核时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='提现申请表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `withdrawals`
--

LOCK TABLES `withdrawals` WRITE;
/*!40000 ALTER TABLE `withdrawals` DISABLE KEYS */;
INSERT INTO `withdrawals` VALUES (1,7,500.00,'工商银行','6222****1234','分销商一','completed','1月提现',2,'2026-01-25 14:00:00','2026-01-24 10:00:00','2026-01-25 14:00:00'),(2,7,300.00,'支付宝','alipay123@dmh.com','分销商一','completed','1月第二笔提现',2,'2026-01-28 16:00:00','2026-01-27 15:00:00','2026-01-28 16:00:00'),(3,8,200.00,'微信支付','wxid123456','分销商二','pending','待审核',NULL,NULL,'2026-01-30 09:00:00','2026-01-30 09:00:00'),(4,2,1000.00,'建设银行','6217****5678','品牌经理','pending','品牌管理员提现',NULL,NULL,'2026-01-31 08:30:00','2026-01-31 08:30:00');
/*!40000 ALTER TABLE `withdrawals` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-01-31 10:48:47
