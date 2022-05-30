package controller

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	action_type := c.Query("action_type")
	dbInit()
	defer db.Close()
	var user []dbUser
	var rootUser []dbUser
	//查询
	db.Select(&user, "select ID, Name, FollowCount, FollowerCount, IsFollow from User where token=?", token)
	//获取粉丝数和关注数全局变量
	db.Select(&rootUser, "select FollowCount, FollowerCount from User where token=?", "rootroooot")

	if user != nil {
		if action_type == "1" {
			db.Exec("update User set IsFollow=? where token=?", token, true)
			//关注数和粉丝数增加为全局变量
			db.Exec("update User set FollowerCount=? where token=?", "rootroooot", rootUser[0].FollowerCount+1)
			db.Exec("update User set FollowCount=? where token=?", "rootroooot", rootUser[0].FollowCount+1)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Follow success"})
		} else if action_type == "2" {
			db.Exec("update User set IsFollow=? where token=?", token, false)
			//关注数和粉丝数为全局变量
			db.Exec("update User set FollowerCount=? where token=?", "rootroooot", rootUser[0].FollowerCount+1)
			db.Exec("update User set FollowCount=? where token=?", "rootroooot", rootUser[0].FollowCount+1)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Unfollow success"})
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}
