package mysql

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

const (
	INSERT_COUNT = 100000
)

// 一条一条插入
func InsertOneByOne(db *gorm.DB) {
	begin := time.Now()
	for i := 0; i < INSERT_COUNT; i++ {
		student := Student{Name: "学生" + strconv.Itoa(i), Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()}
		if err := db.Create(&student).Error; err != nil { //注意需要传student的指针
			fmt.Println(err)
			return
		}
	}
	fmt.Println("total", time.Since(begin))
}

// 放在一个事务里插入
func InsertByTransaction1(db *gorm.DB) {
	begin := time.Now()
	tx := db.Begin()
	for i := 0; i < INSERT_COUNT; i++ {
		student := Student{Name: "学生" + strconv.Itoa(i), Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()}
		if err := tx.Create(&student).Error; err != nil {
			fmt.Println(err)
			return
		}
	}
	tx.Commit()
	fmt.Println("total", time.Since(begin))
}

// 一次插入多条，整体再放到一个事务里
func InsertByTransaction2(db *gorm.DB) {
	begin := time.Now()
	tx := db.Begin()
	const BATCH = 100 // 一条SQL语句插入多条
	for i := 0; i < INSERT_COUNT; i += BATCH {
		students := make([]Student, 0, BATCH)
		for j := 0; j < BATCH; j++ {
			student := Student{Name: "学生" + strconv.Itoa(i+j), Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()}
			students = append(students, student)
		}
		if err := tx.Create(&students).Error; err != nil {
			fmt.Println(err)
			return
		}
	}
	tx.Commit()
	fmt.Println("total", time.Since(begin))
}
