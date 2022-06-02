package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
	//以⼆进制格式上传⽂件
	//这是⼀个Post 参数会被返回的地址
	uri := "http://127.0.0.1:8888/static/bear.mp4"
	byte, err := ioutil.ReadFile("bear.mp4")
	if err != nil {
		fmt.Println("err=", err)
	}
	res, err := http.Post(uri, "multipart/form-data", bytes.NewReader(byte))
	if err != nil {
		fmt.Println("err=", err)
	}
	//http返回的response的body必须close,否则就会有内存泄露
	defer func() {
		res.Body.Close()
	}()
	//读取body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(" post err=", err)
	}
	fmt.Println(string(body))

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: ""},
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}
