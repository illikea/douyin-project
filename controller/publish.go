package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	dbInit()
	defer db.Close()
	var users []dbUser
	//查询
	db.Select(&users, "select ID, Name from User where token=?", token)
	if users == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	title, err := c.FormFile("title")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	filename := filepath.Base(data.Filename)
	//user := usersLoginInfo[token]  默认用户投稿test
	finalName := fmt.Sprintf("%d_%s", users[0].ID, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	db.Exec("insert into Video(ID, Author, PlayUrl, CoverUrl, FavoriteCount, CommentCount, IsFavorite, Title)value(?, ?, ?, ?, ?, ?, ?, ?)", 1, token, "http://127.0.0.1:8080/static/"+filename, "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg", 0, 0, 0, title)
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	token := c.Query("token")
	var videoList []Video
	rows, _ := db.Query("select ID, Author, PlayUrl, CoverUrl, FavoriteCount, CommentCount, IsFavorite, Title from Vedio where ID>?", -1)
	defer rows.Close()
	//获取用户信息
	var user []dbUser
	db.Select(&user, "select ID, Name, FollowCount, FollowerCount, IsFollow from User where token=?", token)
	var ResponseUser = User{
		Id:            user[0].ID,
		Name:          user[0].Name,
		FollowCount:   user[0].FollowCount,
		FollowerCount: user[0].FollowerCount,
		IsFollow:      user[0].IsFollow,
	}
	for rows.Next() {
		var video dbVideo
		rows.Scan(&video.ID, &video.PlayUrl, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &video.IsFavorite, &video.Title)
		videoList = append(videoList, Video{
			Id:            video.ID,
			Author:        ResponseUser,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		VideoList: videoList,
	})
}
