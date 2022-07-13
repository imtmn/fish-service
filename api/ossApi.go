package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
)

func OSSUploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": 502, "msg": "参数有误", "err": err.Error()})
		return
	}
	fileExt := filepath.Ext(fileHeader.Filename)
	allowExts := []string{".jpg", ".png", ".gif", ".jpeg", ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx", ".pdf"}
	allowFlag := false
	for _, ext := range allowExts {
		if ext == fileExt {
			allowFlag = true
			break
		}
	}
	if !allowFlag {
		c.JSON(http.StatusBadGateway, gin.H{"code": 502, "msg": "不允许的类型", "err": err.Error()})
		return
	}

	now := time.Now()
	//文件存放路径
	fileDir := fmt.Sprintf("upload/%s", now.Format("200601"))

	//文件名称
	timeStamp := now.Unix()
	fileName := fmt.Sprintf("%d-%s", timeStamp, fileHeader.Filename)
	// 文件key
	fileKey := filepath.Join(fileDir, fileName)

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": 502, "msg": "上传失败,open oss error", "err": err.Error()})
		return
	}
	defer src.Close()

	res, err := OssUpload(fileKey, src)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": 502, "msg": "OssUpload上传失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "文件上传成功", "path": res})
}

func OssUpload(key string, file io.Reader) (string, error) {
	client, err := oss.New(os.Getenv("OSS_ENDPOINT"), os.Getenv("OSS_ACCESS_KEY_ID"), os.Getenv("OSS_ACCESS_SECRET"))
	if err != nil {
		return "client 创建失败", err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(os.Getenv("OSS_BUCKET"))
	if err != nil {
		return "bucket 创建失败", err
	}
	// 上传文件。
	err = bucket.PutObject(key, file)
	if err != nil {
		return "文件上传失败", err
	}
	return "/" + key, nil
}
