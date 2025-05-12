package mysql

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 默认情况下，GORM 使用 ID 作为主键，使用结构体名的 蛇形复数 作为表名，字段名的 蛇形 作为列名
type Student struct {
	Id         int    `gorm:"column:id;primaryKey"`
	Name       string `gorm:"column:name"`
	Province   string
	City       string    `gorm:"column:city"`
	Address    string    `gorm:"column:addr"`
	Score      float32   `gorm:"column:score"`
	Enrollment time.Time `gorm:"column:enrollment;type:date"`
	Gender     int       `gorm:"-"` //不做映射（比如表里没有这个字段，但结构体需要有它）
}

// 使用TableName()来修改默认的表名
func (Student) TableName() string {
	return "student"
}

const (
	PAGE_SIZE = 100
)

// Traverse1 通过limit offset,n遍历表
func Traverse1(db *gorm.DB) {
	var offset int
	begin := time.Now()
	for i := 0; ; i++ {
		t0 := time.Now()
		var students []Student
		db.Select("id,name").Limit(PAGE_SIZE).Offset(offset).Find(&students)
		offset += len(students)
		if i%300 == 0 {
			fmt.Println(i, time.Since(t0).Microseconds(), "us")
		}
		if len(students) == 0 {
			break
		}
	}
	fmt.Println("total", time.Since(begin))
}

// Traverse2 通过id>x limit 0,n遍历表
func Traverse2(db *gorm.DB) {
	var maxid int
	begin := time.Now()
	for i := 0; ; i++ {
		t0 := time.Now()
		var students []Student
		db.Select("id,name").Where("id>?", maxid).Limit(PAGE_SIZE).Find(&students) //默认会按id排好序
		if i%300 == 0 {
			fmt.Println(i, time.Since(t0).Microseconds(), "us")
		}
		if len(students) == 0 {
			break
		}
		// for _, stu := range students {
		// 	if stu.Id > maxid {
		// 		maxid = stu.Id
		// 	}
		// }
		maxid = students[len(students)-1].Id //已按id排好序，所以最后一个students的Id最大
	}
	fmt.Println("total", time.Since(begin))
}
