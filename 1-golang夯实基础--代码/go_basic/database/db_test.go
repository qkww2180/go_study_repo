package database_test

import (
	"database/sql"
	"dqq/go/basic/database"
	"fmt"
	"testing"
)

var (
	db *sql.DB
)

func init() {
	var err error
	/**
	DSN(data source name)格式：[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	例如user:password@tcp(localhost:5555)/dbname?charset=utf8
	如果是本地MySQl，且采用默认的3306端口，可简写为：user:password@/dbname
	想要正确的处理time.Time ，您需要带上parseTime参数
	要支持完整的UTF-8编码，您需要将charset=utf8更改为charset=utf8mb4
	loc=Local采用机器本地的时区
	*/
	// db, err := sql.Open("mysql", "tester:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	//db可以并发使用
	db, err = sql.Open("mysql", "tester:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
	database.CheckError(err)
	// defer db.Close()
}

func TestInsert(t *testing.T) {
	database.Insert(db)
}

func TestReplace(t *testing.T) {
	database.Replace(db)
}

func TestMassInsert(t *testing.T) {
	database.MassInsertStmt(db)
}

func TestUpdate(t *testing.T) {
	database.Update(db)
}

func TestSelect(t *testing.T) {
	database.Query(db)
}

func TestDelete(t *testing.T) {
	database.Delete(db)
}

func TestTransaction(t *testing.T) {
	database.Transaction(db)
}

func TestQueryByPage(t *testing.T) {
	total, data := database.QueryByPage(db, 10, 2)
	fmt.Println(total)
	for _, user := range data {
		fmt.Println(user.Id, user.Score)
	}
}

func TestSqlBuilder(t *testing.T) {
	database.SqlInsert()
	database.SqlDelete()
	database.SqlRead()
	database.SqlUpdate()
}

// go test -v ./database -run=^TestInsert$ -count=1
// go test -v ./database -run=^TestReplace$ -count=1
// go test -v ./database -run=^TestMassInsert$ -count=1
// go test -v ./database -run=^TestUpdate$ -count=1
// go test -v ./database -run=^TestTransaction$ -count=1
// go test -v ./database -run=^TestSelect$ -count=1
// go test -v ./database -run=^TestDelete$ -count=1
// go test -v ./database -run=^TestSqlBuilder$ -count=1
// go test -v ./database -run=^TestQueryByPage$ -count=1
