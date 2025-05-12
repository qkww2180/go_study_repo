package database

import (
	"context"
	"database/sql"
	"fmt"
	rand "math/rand/v2"
	"time"

	_ "github.com/go-sql-driver/mysql" //init()
	gsb "github.com/huandu/go-sqlbuilder"
)

// insert 插入数据
func Insert(db *sql.DB) {
	//一条sql，插入2行记录
	res, err := db.Exec("insert into student (name,province,city,enrollment) values ('小明', '深圳', '深圳', '2022-07-03'), ('小红', '上海', '上海', '2022-07-03')")
	CheckError(err)
	lastId, err := res.LastInsertId() //ID自增，用过的id（即使对应的行已delete）不会重复使用。如果使用单个INSERT语句将多行插入到表中，则LastInsertId是第一条数据使用的id
	CheckError(err)
	fmt.Printf("after insert last id %d\n", lastId)
	rows, err := res.RowsAffected() //插入2行，所以影响了2行
	CheckError(err)
	fmt.Printf("insert affect %d row\n", rows)
}

// replace 插入(覆盖)数据
func Replace(db *sql.DB) {
	//由于name字段上有唯一索引，insert重复的name会报错。而使用replace会先删除，再插入
	res, err := db.Exec("replace into student (name,province,city,enrollment) values ('小明', '深圳', '深圳', '2025-07-03'), ('小红', '上海', '上海', '2025-07-03')")
	CheckError(err)
	lastId, err := res.LastInsertId() //ID自增，用过的id（即使对应的行已delete）不会重复使用
	CheckError(err)
	fmt.Printf("after insert last id %d\n", lastId)
	rows, err := res.RowsAffected() //先删除，后插入，影响了4行
	CheckError(err)
	fmt.Printf("insert affect %d row\n", rows)
}

// update 修改数据
func Update(db *sql.DB) {
	//不同的city加不同的分数
	res, err := db.Exec("update student set score=score+10 where city='上海'") //上海加10分
	CheckError(err)
	lastId, err := res.LastInsertId() //0, 仅插入操作才会给LastInsertId赋值
	CheckError(err)
	fmt.Printf("after update last id %d\n", lastId)
	rows, err := res.RowsAffected() //where city=?命中了几行，就会影响几行
	CheckError(err)
	fmt.Printf("update affect %d row\n", rows)
}

// 事务
func Transaction(db *sql.DB) {
	tx, err := db.BeginTx(context.Background(), nil) // 开始事务
	CheckError(err)
	_, err = tx.Exec("insert into student (name,province,city,enrollment,score) values ('Tom', '深圳', '深圳', '2022-07-03',40)")
	CheckError(err)
	// _, err = tx.Exec("insert into student (name,province,city,enrollment,score) values ('Tom', '深圳', '深圳', '2022-07-03',40)") // 一旦中间某一步出错失败，则事务里的所有操作全部回滚
	// CheckError(err)
	tx.Commit() //整体提交
}

// delete 删除数据
func Delete(db *sql.DB) {
	res, err := db.Exec("delete from student where id=6") //删除得分大于13的记录
	CheckError(err)
	rows, err := res.RowsAffected() //where id>13命中了几行，就会影响几行
	CheckError(err)
	fmt.Printf("delete affect %d row\n", rows)
}

type User struct {
	Id     int
	Gender string
	Score  float64
}

func QueryByPage(db *sql.DB, pageSize, page int) (total int, data []*User) {
	// 先获得总数
	rows, err := db.Query("select count(*) from student")
	CheckError(err)
	defer rows.Close()
	rows.Next() //一定要先执行Next()
	rows.Scan(&total)

	// 再通过select limit 获得相应分页里的数据
	offset := pageSize * (page - 1)
	rows2, err := db.Query(fmt.Sprintf("select id,score from student limit %d,%d", offset, pageSize))
	CheckError(err)
	defer rows2.Close()
	for rows2.Next() {
		var id int
		var score float64
		rows2.Scan(&id, &score)
		data = append(data, &User{Id: id, Score: score})
	}
	return
}

// query 查询数据
func Query(db *sql.DB) map[int]*User {
	rows, err := db.Query("select id,name,city,score,enrollment from student where enrollment>=20250130 limit 5") //查询得分大于2的记录
	CheckError(err)
	defer rows.Close()
	rect := make(map[int]*User, 10)
	for rows.Next() { //没有数据或发生error时返回false
		var id int
		var score float32
		var name, city string
		var enrollment time.Time
		err = rows.Scan(&id, &name, &city, &score, &enrollment) //通过scan把db里的数据赋给go变量
		CheckError(err)
		fmt.Printf("id=%d, score=%.2f, name=%s, city=%s, enrollment=%s \n", id, score, name, city, enrollment.Format("2006-01-02"))
		rect[id] = &User{
			Id:    id,
			Score: float64(score),
		}
	}
	return rect
}

func QueryUser(db *sql.DB, mp map[int]*User) {
	rows, err := db.Query("select id,gender from user") //查询得分大于2的记录
	CheckError(err)
	defer rows.Close()
	for rows.Next() { //没有数据或发生error时返回false
		var id int
		var gender string
		err = rows.Scan(&id, &gender) //通过scan把db里的数据赋给go变量
		CheckError(err)
		user, ok := mp[id]
		if ok {
			user.Gender = gender
		} else {
			mp[id] = &User{
				Id:     id,
				Gender: gender,
			}
		}
	}
}

// 插入大量记录
func MassInsertStmt(db *sql.DB) {
	insertBuilder := gsb.NewInsertBuilder()
	insertBuilder = insertBuilder.InsertInto("student").Cols("name", "province", "city", "enrollment", "score").Values(RandStringRunes(10), "河南", "郑州", time.Now().Add(time.Hour*24*time.Duration(1)).Format("2006-01-02"), rand.IntN(100))
	sql, args := insertBuilder.Build()
	stmt, err := db.Prepare(sql)
	CheckError(err)
	stmt.Exec(args...)
	for i := 0; i < 1000; i++ {
		stmt.Exec(RandStringRunes(10), "河南", "郑州", time.Now().Add(time.Hour*24*time.Duration(i)).Format("2006-01-02"), rand.IntN(100)) //重复利用stmt
	}
	stmt.Close()

	sb := gsb.NewInsertBuilder()
	sb = sb.InsertInto("student").Cols("name", "province", "city", "enrollment", "score")
	for i := 0; i < 1000; i++ {
		sb = sb.Values(RandStringRunes(10), "河南", "郑州", time.Now().Add(time.Hour*24*time.Duration(i)).Format("2006-01-02"), rand.IntN(100))
	}
	sql, args = sb.Build()
	stmt2, err := db.Prepare(sql)
	CheckError(err)
	stmt2.Exec(args...)
	stmt2.Close()
}
