package mysql

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestInsertOneByOne(t *testing.T) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		panic(err)
	}
	db.Delete(Student{}, "1=1") //把表清空

	InsertOneByOne(db)
}

func TestInsertByTransaction1(t *testing.T) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		panic(err)
	}
	db.Delete(Student{}, "1=1") //把表清空

	InsertByTransaction1(db)
}

func TestInsertByTransaction2(t *testing.T) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		panic(err)
	}
	db.Delete(Student{}, "1=1") //把表清空

	InsertByTransaction2(db)
}

// go test -v ./database/mysql -run=^TestInsert -count=1

/**
=== RUN   TestInsertOneByOne
total 3.584101s
--- PASS: TestInsertOneByOne (3.60s)
=== RUN   TestInsertByTransaction1
total 212.886ms
--- PASS: TestInsertByTransaction1 (0.22s)
=== RUN   TestInsertByTransaction2
total 46.7347ms
--- PASS: TestInsertByTransaction2 (0.06s)
PASS
*/
