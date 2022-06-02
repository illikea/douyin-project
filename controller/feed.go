package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const BaseUploadPath = "/public/"

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	/*//以⼆进制格式上传⽂件
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
	fmt.Println(string(body))*/

	http.HandleFunc("/download", handleDownload)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("Server run fail")
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: ""},
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}

//文件下载
func handleDownload(w http.ResponseWriter, request *http.Request) {
	//文件上传只允许GET方法
	if request.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method not allowed"))
		return
	}
	//文件名
	filename := request.FormValue("filename")
	if filename == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
	log.Println("filename: " + filename)
	//打开文件
	file, err := os.Open(BaseUploadPath + filename)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
	//结束后关闭文件
	defer file.Close()

	//设置响应的header头
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
	//将文件写至responseBody
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
}
