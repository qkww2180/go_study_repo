package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bytedance/sonic"
	_ "github.com/mailru/go-clickhouse/v2"
)

var (
	Loc *time.Location
)

func init() {
	Loc, _ = time.LoadLocation("Asia/Shanghai")
}

type userInfo struct {
	Phone  string `json:"phone"`  //手机操作系统
	Gender byte   `json:"gender"` //1男，2女
}

// 批量导入某天注册的用户数据（通常是每天凌晨，从MySQL向ClickHouse导一批数据）
func batchImport(conn *sql.DB, day time.Time) {
	day = time.Unix(86400*day.Unix()/86400, 0) //归到0点
	dayStr := day.Format("2016-01-02")

	tx, err := conn.Begin() //开始事务
	if err != nil {
		log.Printf("begin transaction failed: %v", err)
		return
	}

	//导入之前，先把当天已有的数据全部删除掉。确保batchImport是个幂等函数
	deleteSQL := "delete from user where date(create_time)='" + dayStr + "'"
	if _, err := tx.Exec(deleteSQL); err != nil {
		log.Printf("删除%s的user表数据失败:%v", dayStr, err)
		return
	}

	stmt, err := tx.Prepare(`insert into user (user_id,name,create_time,extra) values (?,?,?,?)`)
	if err != nil {
		log.Printf("prepare statement failed: %v", err)
		return
	}

	//模拟大量假数据，插入CH
	total := 2000 + rand.Intn(3000) //total是2000到5000上的随机数
	for i := 0; i <= total; i++ {
		userId := rand.Intn(99999999)
		name := "大乔乔"
		createTime := time.Unix(day.Unix()+int64(rand.Intn(86400)), 0)
		info := userInfo{}
		if rand.Int()%2 == 0 {
			info.Phone = "ios"
			info.Gender = 1
		} else {
			info.Phone = "android"
			info.Gender = 2
		}
		if infoBytes, err := sonic.Marshal(info); err == nil {
			if _, err = stmt.Exec(
				userId,
				name,
				createTime,
				string(infoBytes), //在CH里是JSON类型，insert时用string(json序列化之后)类型
			); err != nil {
				log.Printf("插入user数据失败:%v", err)
			}
		} else {
			log.Printf("json序列化出错:%v", err)
		}
	}
	if err = tx.Commit(); err != nil { //提交事务
		log.Printf("commit batch insert failed: %v", err)
	}
}

// 查询一段时间[begin, end)内每天的用户注册量
func query(conn *sql.DB, begin, end time.Time) []int {
	rect := make([]int, 0, 10)
	sql := "select date(create_time) as date,uniq(user_id) from user where date>='" + begin.Format("2016-01-02") + "' and date<'" + end.Format("2016-01-02") + "' group by date order by date"
	if rows, err := conn.Query(sql); err == nil {
		for rows.Next() {
			var day time.Time
			var count int
			if err := rows.Scan(&day, &count); err == nil {
				rect = append(rect, count)
			} else {
				log.Printf("scan user failed: %v", err)
			}
		}
	} else {
		log.Printf("query table user failed: %v", err)
	}
	return rect
}

func main() {
	host := "localhost"                                                                                                                    //连接本机的ClickHouse
	port := 8123                                                                                                                           //ClickHouse Server默认使用的端口
	db := "test"                                                                                                                           //使用哪个数据库
	user := "default"                                                                                                                      //ClickHouse默认会创建一个叫"default"的用户
	passwd := ""                                                                                                                           //删除ubuntu上的文件/etc/clickhouse-server/users.d/default-password.xml，default用户登录就不需要密码了
	conn, err := sql.Open("chhttp", fmt.Sprintf("http://%s:%s@%s:%d/%s?read_timeout=10s&write_timeout=20s", user, passwd, host, port, db)) //通过http协议连接
	if err != nil {
		panic(err)
	}
	if err := conn.Ping(); err != nil {
		panic(err)
	}
	defer conn.Close() //用完后，关闭数据库连接

	//往CH的user表写入7天的数据
	day, _ := time.ParseInLocation("2006-01-02", "2023-09-01", Loc)
	for i := 0; i < 7; i++ {
		batchImport(conn, day)
		day = day.Add(24 * time.Hour)
	}

	begin, _ := time.ParseInLocation("2006-01-02", "2023-09-01", Loc)
	end, _ := time.ParseInLocation("2006-01-02", "2023-09-08", Loc)
	counts := query(conn, begin, end)
	fmt.Println(counts)
}

// go run .\g_database\clickhouse\
