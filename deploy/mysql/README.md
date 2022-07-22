# 配置主从（docker创建主从模式的数据库实例）

## 1、开启一个主节点、一个从节点  
```
docker run --name mysql-master -p 3306:3306 -e MYSQL_ROOT_PASSWORD=12345678 -v E:/Desktop/项目/cloudStorage/deploy/mysql/master/data:/var/lib/mysql -d mysql:latest #-v data目录映射
docker run --name mysql-slave -p 3307:3306 -e MYSQL_ROOT_PASSWORD=12345678 -d mysql:latest
```
![图片(./res/20220715235445.png)

## 2、修改数据库server-id，防止主从模式的数据库实例冲突  
```
#拷贝配置文件
docker cp mysql-master:/etc/my.cnf ./master
docker cp mysql-slave:/etc/my.cnf ./slave01

#修改配置文件，更改server-id，主节点server-id为1000，从节点server-id为1001
......

#覆盖配置文件
docker cp ./master/my.cnf mysql-master:/etc/my.cnf
docker cp ./slave01/my.cnf mysql-slave:/etc/my.cnf
```

## 3、docker重启数据库实例  
```
docker restart mysql-master
docker restart mysql-slave

#查看server-id
show VARIABLES like 'server_id';
```

## 4、主节点查看状态  
```
show master status; #可以看到binlog文件的位置
```
binlog:binlog.000002  

## 5、从节点指定主节点的信息  
```
CHANGE MASTER TO MASTER_HOST = '主节点的IP地址(不能使用127.0.0.1/localhost)', MASTER_PORT = 3306, MASTER_USER = 'root', MASTER_PASSWORD = '密码', MASTER_LOG_FILE = 'mysql-bin.000002', MASTER_LOG_POS = 0;

start slave; # 启动从节点
show slave status\G; # 查看从节点的状态
```

## 6、主节点创建数据库表  
```
#创建数据库
create database test;

#创建测试表
CREATE TABLE `test` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

## 7、开启远程访问  
```
GRANT ALL PRIVILEGES ON *.* TO  'root'@'%' with grant option;
flush privileges;
```

# 数据表
```
#tbl_file
create database tbl_file;
use tbl_file;
CREATE TABLE `tbl_file` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `file_shal` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
    `file_size` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件大小',
    `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
    `create_at` datetime DEFAULT NOW() COMMENT '创建时间',
    `update_at` datetime DEFAULT NOW() on update current_timestamp() COMMENT '更新时间',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等)',
    `ext1` int(11) DEFAULT '0' COMMENT '扩展字段1',
    `ext2` text COMMENT '扩展字段2',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_file_shal` (`file_shal`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件表';

#tbl_user
CREATE TABLE `tbl_user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
    `user_pwd` varchar(64) NOT NULL DEFAULT '' COMMENT '用户密码',
    `email` varchar(64) DEFAULT '' COMMENT '邮箱',
    `phone` varchar(128) DEFAULT '' COMMENT '手机号',
    `email_vaildated` tinyint(1) DEFAULT 0 COMMENT '邮箱是否验证',
    `phone_vaildated` tinyint(1) DEFAULT 0 COMMENT '手机号是否验证',
    `signup_at` datetime DEFAULT current_timestamp() COMMENT '注册时间',
    `update_at` datetime DEFAULT current_timestamp() on update current_timestamp() COMMENT '最后活跃时间',
    `profile` text COMMENT '用户属性',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_phone` (`phone`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

#tbl_user_token
CREATE TABLE `tbl_user_token` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
    `user_token` char(40) NOT NULL DEFAULT '' COMMENT 'token',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_name` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户token表';

#tbl_user_file
CREATE TABLE `tbl_user_file` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(64) NOT NULL COMMENT '用户名',
    `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
    `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
    `update_at` datetime DEFAULT current_timestamp() on update current_timestamp() COMMENT '最后修改时间',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常/1删除/2禁用)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_file` (`file_name`,`file_sha1`),
    KEY `idx_user_name` (`user_name`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户文件表';


```

