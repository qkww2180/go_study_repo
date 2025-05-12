package mysql

import (
	"gorm.io/gorm"
)

// 登录成功返回true。容易被SQL注入攻击
func LoginUnsafe(db *gorm.DB, name, passwd string) bool {
	var cnt int64
	db.Table("login").Select("*").Where("username='" + name + "' and password='" + passwd + "'").Count(&cnt)
	return cnt > 0
}

// 登录成功返回true。拒绝SQL注入攻击
func LoginSafe(db *gorm.DB, name, passwd string) bool {
	var cnt int64
	db.Table("login").Select("*").Where("username=? and password=?", name, passwd).Count(&cnt)
	return cnt > 0
}