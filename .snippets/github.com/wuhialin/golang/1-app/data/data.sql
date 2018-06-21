CREATE TABLE `wowpower_game` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '价格',
  `day` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '周期',
  `gain_profit` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '返利',
  `create_date` date NOT NULL COMMENT '日期',
  `update_time` int(10) unsigned NOT NULL COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `create_date` (`create_date`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='蛙宝游戏';