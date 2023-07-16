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