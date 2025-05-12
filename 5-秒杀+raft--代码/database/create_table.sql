CREATE TABLE `inventory` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '奖品id，自增',
  `name` varchar(20) NOT NULL COMMENT '奖品名称',
  `description` varchar(100) NOT NULL DEFAULT '' COMMENT '奖品描述',
  `picture` varchar(200) NOT NULL DEFAULT '' COMMENT '奖品图片',
  `price` int(11) NOT NULL DEFAULT '0' COMMENT '价值',
  `count` int(11) NOT NULL DEFAULT '0' COMMENT '库存量',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8 COMMENT='奖品库存表';

CREATE TABLE `orders` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '订单id，自增',
  `gift_id` int(11) NOT NULL COMMENT '商品id',
  `user_id` int(11) NOT NULL COMMENT '用户id',
  `count` int(11) NOT NULL DEFAULT '1' COMMENT '购买数量',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '订单创建时间',       
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11011 DEFAULT CHARSET=utf8mb4 COMMENT='订单表';