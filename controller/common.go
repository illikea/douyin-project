package controller

import (
	"fmt"
	idworker "github.com/gitstliu/go-id-worker"
	"github.com/jmoiron/sqlx"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

var db *sqlx.DB

func dbInit() {
	database, err := sqlx.Open("mysql", "root:984435589dsaqY@tcp(localhost:3306)/douyint") //flash123为MySQL密码，需改
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	db = database
	//因为会提前关闭，暂时不关闭
	//defer db.Close() // 注意这行代码要写在上面err判断的下面
}

type dbUser struct {
	ID            int64  `db:"ID"`
	Name          string `db:"Name"`
	FollowCount   int64  `db:"FollowCount"`
	FollowerCount int64  `db:"FollowerCount"`
	IsFollow      bool   `db:"IsFollow"`
	token         string `db:"token"`
}

type dbVideo struct {
	ID            int64  `db:"ID"`
	Author        string `db:"Author"`
	PlayUrl       string `db:"PlayUrl"`
	CoverUrl      string `db:"CoverUrl"`
	FavoriteCount int64  `db:"FavoriteCount"`
	CommentCount  int64  `db:"CommentCount"`
	IsFavorite    bool   `db:"IsFavorite"`
	Title         string `db:"Title"`
}

//生成唯一ID
func makeId() int64 {
	currWoker := &idworker.IdWorker{}
	currWoker.InitIdWorker(1000, 1)
	newID, _ := currWoker.NextId()
	return newID
}
