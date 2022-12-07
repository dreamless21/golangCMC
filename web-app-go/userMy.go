package web_app_go

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id       int
	Username string
	Balance  float64
}

func OpenConnection() *sql.DB {
	src := fmt.Sprintf("gomysql:%s@tcp(%s)/golang_webapp", PassMysql, IpMysql)
	db, _ := sql.Open("mysql", src)
	return db
}

func NewUser(username string, balance float64) User {
	db := OpenConnection()
	defer db.Close()
	sqlStmt, _ := db.Prepare("INSERT INTO web_users(username, balance) VALUES(?, ?)")
	res, _ := sqlStmt.Exec(username, balance)
	lastId, _ := res.LastInsertId()
	return User{int(lastId), username, balance}
}

func (u *User) ChangeUsername(newName string) {
	db := OpenConnection()
	defer db.Close()
	sql, _ := db.Prepare("UPDATE web_users SET username = ? WHERE id = ?")
	sql.Exec(newName, u.Id)
	u.Username = newName
}

func GetUser(id int) User {
	var us User
	db := OpenConnection()
	defer db.Close()
	sql := db.QueryRow("SELECT id, username, balance FROM web_users WHERE id = ?", id)
	_ = sql.Scan(&us.Id, &us.Username, &us.Balance)
	return us
}
