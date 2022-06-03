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
	actionType := c.Query("action_type")
	toUserID := c.Query("to_user_id")
	dbInit()
	defer db.Close()
	var user []dbUser
	var toUser []dbUser
	//查询用户是否存在
	db.Select(&user, "select ID, Name, FollowCount, FollowerCount, IsFollow from User where token=?", token)
	db.Select(&toUser, "select ID, Name, FollowCount, FollowerCount, IsFollow, token from User where ID=?", toUserID)

	if user != nil && toUser != nil {
		if actionType == "1" {
			db.Exec("update User set IsFollow=? where token=?", token, true)
			//修改用户关注数和粉丝数，并在FollowList新增一行
			db.Exec("update User set FollowerCount=? where ID=?", toUser[0].ID, toUser[0].FollowerCount+1)
			db.Exec("update User set FollowCount=? where ID=?", user[0].ID, user[0].FollowCount+1)
			db.Exec("insert into FollowList(FollowCount, FollowerCount, FollowerID, UserID, IsFollow, Name, token)value(?, ?, ?, ?, ?, ?, ?)", 0, 0, user[0].ID, toUserID, 1, toUser[0].Name, toUser[0].token)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Follow success"})
		} else if actionType == "2" {
			db.Exec("update User set IsFollow=? where token=?", token, false)
			//修改用户关注数和粉丝数，并在FollowList删除对应行
			db.Exec("update User set FollowerCount=? where ID=?", toUser[0].ID, toUser[0].FollowerCount-1)
			db.Exec("update User set FollowCount=? where ID=?", user[0].ID, user[0].FollowCount-1)
			db.Exec("delete from FollowList where UserID=? AND FollowerID=?", toUser[0].ID, user[0].ID)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Unfollow success"})
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userID := c.Query("user_id")
	dbInit()
	defer db.Close()
	var userList []User
	//从数据库获取关注列表
	var followList []dbFollower
	db.Select(&followList, "select FollowerID, Name, FollowCount, FollowerCount, IsFollow from FollowList where UserID=?", userID)
	//填充至返回的列表
	for _, user := range followList {
		userList = append(userList, User{
			Id:            user.FollowerID,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      user.IsFollow,
		})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: userList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	userID := c.Query("user_id")
	dbInit()
	defer db.Close()
	var userList []User
	//从数据库获取粉丝列表
	var followList []dbFollower
	db.Select(&followList, "select UserID, Name, FollowCount, FollowerCount, IsFollow from FollowList where FollowerID=?", userID)
	//填充至返回的列表
	for _, user := range followList {
		userList = append(userList, User{
			Id:            user.UserID,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      user.IsFollow,
		})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: userList,
	})
}
