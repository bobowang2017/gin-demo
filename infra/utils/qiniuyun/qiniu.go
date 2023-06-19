package qiniuyun

import (
	"fmt"
	"gin-demo/core/settings"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"sync"
)

var (
	qiNiuClient *QiNiuClient
	qiNiuOnce   sync.Once
)

type QiNiuClient struct {
	accessKey string
	secretKey string
}

func GetQiNiuClient() *QiNiuClient {
	qiNiuOnce.Do(func() {
		qiNiuClient = &QiNiuClient{
			accessKey: settings.Config.QiNiuYun.AccessKey,
			secretKey: settings.Config.QiNiuYun.SecretKey,
		}
	})
	return qiNiuClient
}

// GetUploadToken 获取上传凭证
func (q *QiNiuClient) GetUploadToken(bucket string, expire uint64) string {
	putPolicy := storage.PutPolicy{
		Scope:   bucket,
		Expires: expire, //示例2小时有效期(7200)
	}
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	return putPolicy.UploadToken(mac)
}

// RefreshUploadToken 覆盖上传的凭证,覆盖过程中需要指定文件名称
func (q *QiNiuClient) RefreshUploadToken(bucket, fileName string) string {
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", bucket, fileName),
	}
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	return putPolicy.UploadToken(mac)
}

// SDK文档 https://developer.qiniu.com/kodo/1238/go
