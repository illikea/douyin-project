package controller

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"sync/atomic"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	/*"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},*/
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	newID := makeId()
	token := username + password

	dbInit()
	defer db.Close()
	var user []dbUser
	//查询,保证用户名唯一
	db.Select(&user, "select ID from User where Name=?", username)
	//若查询到则直接返回
	if user != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, newID)
		_, err := db.Exec("insert into User(FollowCount, FollowerCount, ID, IsFollow, Name, token)value(?, ?, ?, ?, ?, ?)", 0, 0, newID, 0, username, token)
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User register fail"},
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0, StatusMsg: "User register success"},
				UserId:   newID,
				Token:    username + password,
			})
		}
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	dbInit()
	defer db.Close()
	var user []dbUser
	//查询
	db.Select(&user, "select ID from User where token=?", token)
	if user != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user[0].ID,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist or password error"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	dbInit()
	defer db.Close()
	var user []dbUser
	//查询
	db.Select(&user, "select ID, Name, FollowCount, FollowerCount, IsFollow from User where token=?", token)
	if user != nil {
		var ResponseUser = User{
			Id:            user[0].ID,
			Name:          user[0].Name,
			FollowCount:   user[0].FollowCount,
			FollowerCount: user[0].FollowerCount,
			IsFollow:      user[0].IsFollow,
		}
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0, StatusMsg: ""},
			User:     ResponseUser,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
