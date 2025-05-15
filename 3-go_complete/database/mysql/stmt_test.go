package mysql

import (
	"log"
	"os"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	ormlog "gorm.j_io/gorm/logger"
)

var (
	dsn          = "tester:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	name, passwd = "tom", "123456"
)

func TestLoginUnsafe(t *testing.T) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true, Logger: ormlog.New(log.New(os.Stdout, "\r\n", log.LstdFlags), ormlog.Config{
		SlowThreshold: 100 * time.Millisecond, // 慢 SQL 阈值
		LogLevel:      ormlog.Info,            // Log level，Silent表示不输出日志
		Colorful:      true,                   // 禁用彩色打印
	})}) //强行使用PrepareStmt
	if err != nil {
		panic(err)
	}
	if LoginUnsafe(db, "tom", "123456") == false {
		t.Fail()
	}
	if LoginUnsafe(db, "tom", "456789") == true {
		t.Fail()
	}
	// select * from login where username='tom' and password='456789' or '1'='1'
	if LoginUnsafe(db, "tom", "456789' or '1'='1") == false {
		t.Fail()
	}
}

func TestLoginSafe(t *testing.T) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true, Logger: ormlog.New(log.New(os.Stdout, "\r\n", log.LstdFlags), ormlog.Config{
		SlowThreshold: 100 * time.Millisecond, // 慢 SQL 阈值
		LogLevel:      ormlog.Info,            // Log level，Silent表示不输出日志
		Colorful:      true,                   // 禁用彩色打印
	})}) //强行使用PrepareStmt
	if err != nil {
		panic(err)
	}
	if LoginSafe(db, "tom", "123456") == false {
		t.Fail()
	}
	if LoginSafe(db, "tom", "456789") == true {
		t.Fail()
	}
	if LoginSafe(db, "tom", "456789' or '1'='1") == true {
		t.Fail()
	}
}

func BenchmarkQueryWithoutPrepare(b *testing.B) {
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) //没有指定PrepareStmt
	if err != nil {
		panic(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		LoginUnsafe(client, name, passwd)
	}
}

func BenchmarkQueryWithPrepare(b *testing.B) {
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true}) //强行使用PrepareStmt
	if err != nil {
		panic(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		LoginUnsafe(client, name, passwd)
	}
}

// go test -v ./g_database/mysql -run=^TestLoginSafe -count=1
// go test ./g_database/mysql -bench=^BenchmarkQueryWith.*Prepare$ -run=^$ -count=1 -benchmem
