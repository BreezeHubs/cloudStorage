package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

var cephConn *s3.S3

func GetCephConnection() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}

	//初始化ceph的信息
	auth := aws.Auth{
		AccessKey: "4R6VUG9BLUVTDAZIQ231",
		SecretKey: "yRHsmkfKaedtUxfGKJdULNVAfNG9q5mO9lxmb4ZP",
	}

	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          "http://127.0.0.1:9080", //web服务的地址
		S3Endpoint:           "http://127.0.0.1:9080", //s3服务的地址
		S3BucketEndpoint:     "",                      //s3的bucket的地址
		S3LocationConstraint: false,                   //true表示使用国际化endpoint，false表示使用中国区域endpoint
		S3LowercaseBucket:    false,                   //创建bucket时是否转换为小写
		Sign:                 aws.SignV2,              //签名方式
	}

	//创建ceph的连接
	return s3.New(auth, curRegion)
}

//获取ceph的bucket
func GetCephBucket(bucketName string) *s3.Bucket {
	return GetCephConnection().Bucket(bucketName)
}
