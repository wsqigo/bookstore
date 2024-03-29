CREATE TABLE IF NOT EXISTS `user` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL UNIQUE COMMENT '用户编号',
    `username` varchar(64) NOT NULL UNIQUE COMMENT '用户名',
    `password` varchar(64) NOT NULL COMMENT '密码',
    `email` varchar(64) NOT NULL UNIQUE COMMENT '邮箱',
    `gender` tinyint(4) NOT NULL DEFAULT '0',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

use sql_demo;

DROP TABLE if exists `community`;
CREATE TABLE `community` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `community_id` int(10) unsigned NOT NULL,
    `community_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
    `introduction` varchar(256) COLLATE utf8mb4_general_ci NOT NULL,
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `community_id_idx` (`community_id`),
    UNIQUE KEY  `community_name_idx` (`community_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `community` VALUES ('1', '1', 'Go', 'Golang', '2016-11-01 08:10:10', '2016-11-01 08:10:10');
INSERT INTO `community` VALUES ('2', '2', 'leetcode', '刷题刷题刷题', '2020-01-01 08:00:00', '2020-01-01 08:00:00');
INSERT INTO `community` VALUES ('3', '3', 'CS:GO', 'Rush B...', '2018-08-07 08:30:00', '2018-08-07 08:30:00');
INSERT INTO `community` VALUES ('4', '4', 'LOL', '欢迎来到英雄联盟', '2016-01-01 08:00:00', '2016-01-01 08:00:00');

DROP TABLE IF EXISTS `post`;

CREATE TABLE `post` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `post_id` bigint(20) NOT NULL COMMENT '帖子id',
    `title` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
    `content` varchar(8192) COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
    `author_id` bigint(20) NOT NULL COMMENT '作者的用户id',
    `community_id` bigint(20) NOT NULL COMMENT '所属社区',
    `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '帖子状态',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    primary key (`id`),
    UNIQUE KEY `post_id_idx` (`post_id`),
    KEY `author_id_idx` (`author_id`),
    KEY `community_id_idx` (`community_id`)
) ENGINE=Innodb DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;