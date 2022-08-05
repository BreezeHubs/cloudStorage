package oss

import (
	"cloudStorage/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
)

var ossCli *oss.Client

func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}

	ossCli, err := oss.New(config.OSS_ENDPOINT, config.OSS_ACCESS_KEY, config.OSS_ACCESS_SECRET)
	if err != nil {
		panic(err)
	}
	return ossCli
}

//获取bucket
func Bucket() (*oss.Bucket, error) {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(config.OSS_BUCKET_NAME)
		if err != nil {
			return nil, err
		}
		return bucket, nil
	}
	return nil, errors.New("oss client is nil")
}

//临时授权下载文件
func DownloadUrl(objectName string) (string, error) {
	b, err := Bucket()
	if err != nil {
		return "", errors.Wrap(err, "get bucket error")
	}
	signedURL, err := b.SignURL(objectName, oss.HTTPGet, 3600)
	if err != nil {
		return "", errors.Wrap(err, "SignURL error")
	}
	return signedURL, nil
}

//指定bucket设置生命周期规则
func BuildLifecycleRule(bucketName string) {
	//表示前缀为test的对象（文件）距最后修改时间30天后过期
	ruleTest1 := oss.BuildLifecycleRuleByDays("rule1", "test/", true, 30)
	rules := []oss.LifecycleRule{ruleTest1}

	Client().SetBucketLifecycle(bucketName, rules)
}
