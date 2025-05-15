package mysql

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestTraverse1(t *testing.T) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		panic(err)
	}

	Traverse1(db)
}

func TestTraverse2(t *testing.T) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		panic(err)
	}

	Traverse2(db)
}

// go test -v ./g_database/mysql -run=^TestTraverse -count=1

/**
=== RUN   TestTraverse1
0 1693 us
300 4536 us
600 8341 us
900 12389 us
total 6.9647175s
--- PASS: TestTraverse1 (6.97s)
=== RUN   TestTraverse2
0 4864 us
300 167 us
600 0 us
900 0 us
total 322.7406ms
--- PASS: TestTraverse2 (0.32s)
PASS
*/
