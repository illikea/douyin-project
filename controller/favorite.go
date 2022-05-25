package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	action_type := c.Query("action_type")

	if _, exist := usersLoginInfo[token]; exist {
		if action_type == "1" {
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Like success"})
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Unlike success"})
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
