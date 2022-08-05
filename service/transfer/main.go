package main

import (
	"bufio"
	"cloudStorage/config"
	"cloudStorage/dao"
	"cloudStorage/dao/mq"
	"cloudStorage/store/oss"
	"encoding/json"
	"log"
	"os"
)

func main() {
	log.Println("开始监听转移任务队列...")
	mq.StartConsumer(
		config.TRANS_OSS_QUEUE_NAME,
		"transfer_oss",
		ProcessTransfer,
	)
}

func ProcessTransfer(msg []byte) bool {
	//解析msg
	var pubData mq.TransferData
	if err := json.Unmarshal(msg, &pubData); err != nil {
		log.Printf("解析消息失败，err:%v\n", err)
		return false
	}

	//根据临时存储文件的路径，创建文件句柄
	filed, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Printf("打开文件失败，err:%v\n", err)
		return false
	}

	//通过文件句柄将文件内容读出来并且上传到oss
	b, err := oss.Bucket()
	if err != nil {
		log.Printf("获取bucket失败，err:%v\n", err)
		return false
	}
	if err := b.PutObject(
		pubData.DestLocation,
		bufio.NewReader(filed),
	); err != nil {
		log.Printf("上传文件失败，err:%v\n", err)
		return false
	}

	//更新文件的路径到文件表
	if err := dao.UpdateFileLocation(pubData.FileHash, pubData.DestLocation); err != nil {
		log.Printf("更新数据库的文件路径失败，err:%v\n", err)
		return false
	}

	return true
}
