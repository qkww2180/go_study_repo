CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键自增id',
  `uid` int(11) NOT NULL COMMENT '用户id',
  `keywords` text NOT NULL COMMENT '索引词',
  `degree` char(2) NOT NULL COMMENT '学历',
  `gender` char(1) NOT NULL COMMENT '性别',
  `city` char(2) NOT NULL COMMENT '城市',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid` (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=1707919 DEFAULT CHARSET=utf8 COMMENT='用户信息表';

CREATE TABLE `student` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键自增id',
  `name` char(10) NOT NULL COMMENT '姓名',
  `province` char(6) NOT NULL COMMENT '省',
  `city` char(10) NOT NULL COMMENT '城市',
  `addr` varchar(100) DEFAULT '' COMMENT '地址',
  `score` float NOT NULL DEFAULT '0' COMMENT '考试成绩',
  `enrollment` date NOT NULL COMMENT '入学时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_location` (`province`,`city`)
) ENGINE=InnoDB AUTO_INCREMENT=200001 DEFAULT CHARSET=utf8 COMMENT='学员基本信息';

CREATE TABLE `login` (
  `username` varchar(100) DEFAULT NULL,
  `password` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;