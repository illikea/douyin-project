package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type dbUser struct {
	ID            int64  `db:"ID"`
	Name          string `db:"Name"`
	FollowCount   int64  `db:"FollowCount"`
	FollowerCount int64  `db:"FollowerCount"`
	IsFollow      bool   `db:"IsFollow"`
	token         string `db:"token"`
}

func dbInit() {
	database, err := sqlx.Open("mysql", "root:flash123@tcp(localhost:3306)/douyint") //flash123为MySQL密码，需改
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	db = database
}

func test() {
	username := "12345"
	token := "12345123456"
	dbInit()
	defer db.Close()
	var user []dbUser
	//查询
	err := db.Select(&user, "select ID from User where token=?", token)
	//若查询到则直接返回
	if user != nil {
		fmt.Println("existed\n", err)
		return
	} else {
		//没查到则插入
		db.Exec("insert into User(FollowCount, FollowerCount, ID, IsFollow, Name, token)value(?, ?, ?, ?, ?, ?)", 0, 0, 0, 0, username, token)
	}
}

func main() {
	test()
}
