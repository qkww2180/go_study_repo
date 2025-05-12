package orm

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/gorm"
)

type User struct {
	Id         int    `gorm:"column:id;primaryKey"`
	Gender     string `gorm:"column:sex"` //需要做映射的Field必须是可导出的，否则从DB里查询出来不能给结构体赋值
	Name       string
	FamilyName string `gorm:"-"` //family_name
}

func (User) TableName() string {
	return "user"
}

var (
	_user_fields = GetGormFields(User{})
)

func GetUserById(db *gorm.DB, id int) *User {
	var user User
	if err := db.Select(_user_fields).Where("id=?", id).First(&user).Error; err != nil {
		log.Printf("GetUserById %d failed: %v", id, err)
		return nil
	} else {
		user.FamilyName = user.Name[:3]
		return &user
	}
}

func GetUserByName(db *gorm.DB, name string) *User {
	var user User
	if err := db.Select(_user_fields).Where("name=?", name).First(&user).Error; err != nil {
		log.Printf("GetUserByName %s failed: %v", name, err)
		return nil
	} else {
		user.FamilyName = user.Name[:3]
		return &user
	}
}

func TestGetGormFields(t *testing.T) {
	cols := GetGormFields(User{})
	fmt.Println(cols)
}

// go test -v ./database/orm/ -run=^TestGetGormFields$ -count=1
