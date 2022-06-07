package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	token := c.Query("token")

	dbInit()
	defer db.Close()
	var userLogin []dbUser
	//查询登录用户信息
	db.Select(&userLogin, "select ID, Name, FollowCount, FollowerCount, IsFollow from User where token=?", token)
	var videoList []Video
	//获取视频列表
	rows, _ := db.Query("select ID, AuthorID, PlayUrl, CoverUrl, FavoriteCount, CommentCount, IsFavorite, Title from Video where ID>?", 0)
	//填充视频列表
	if rows != nil {
		for rows.Next() {
			var video dbVideo
			rows.Scan(&video.ID, &video.AuthorID, &video.PlayUrl, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &video.IsFavorite, &video.Title)
			//获取用户信息
			var users []dbUser
			db.Select(&users, "select ID, Name, FollowCount, FollowerCount, IsFollow from User where ID=?", video.AuthorID)
			//若用户已登录，则判断是否已关注视频作者，否则默认未关注
			if userLogin != nil {
				var follow []dbFollower
				db.Select(&follow, "select IsFollow from FollowList where FollowerID=? and UserID=?", userLogin[0].ID, users[0].ID)
				if follow != nil {
					users[0].IsFollow = true
				} else {
					users[0].IsFollow = false
				}
			} else {
				users[0].IsFollow = false
			}
			var user = User{
				Id:            users[0].ID,
				Name:          users[0].Name,
				FollowCount:   users[0].FollowCount,
				FollowerCount: users[0].FollowerCount,
				IsFollow:      users[0].IsFollow,
			}
			videoList = append([]Video{
				{
					Id:            video.ID,
					Author:        user,
					PlayUrl:       video.PlayUrl,
					CoverUrl:      video.CoverUrl,
					FavoriteCount: video.FavoriteCount,
					CommentCount:  video.CommentCount,
					IsFavorite:    video.IsFavorite,
				},
			}, videoList...)
		}
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: ""},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
}
