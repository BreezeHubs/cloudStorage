package main

import (
	"cloudStorage/store/ceph"
	"fmt"
	"os"

	"gopkg.in/amz.v1/s3"
)

func main() {
	bucket := ceph.GetCephBucket("userfile")

	// setBucket(bucket)
	getBucket(bucket)
}

func getBucket(bucket *s3.Bucket) {
	d, _ := bucket.Get("/ceph/b697aa691e55ab94d33be3af5bc61cb47ea71795.docx")
	tmpFile, _ := os.Create("./b697aa691e55ab94d33be3af5bc61cb47ea71795.docx")
	tmpFile.Write(d)
	tmpFile.Close()
}

func setBucket(bucket *s3.Bucket) {
	//创建一个新的bucket
	err := bucket.PutBucket(s3.PublicRead)
	fmt.Println("创建bucket的结果：", err)

	//查询这个bucket下面指定条件的object keys
	resp, err := bucket.List("", "", "", 100)
	fmt.Println("查询bucket下面的object keys的结果：", resp, err)

	//新上传一个对象
	err = bucket.Put("./ceph/test.txt", []byte("hello world"), "octet-stream", s3.PublicRead)
	fmt.Println("新上传一个对象的结果：", err)

	//查询这个bucket下面指定条件的object keys
	resp, err = bucket.List("", "", "", 100)
	fmt.Println("查询bucket下面的object keys的结果：", resp, err)
}
