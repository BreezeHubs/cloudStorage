package handler

import (
	"cloudStorage/config"
	"cloudStorage/dao"
	myredis "cloudStorage/dao/redis"
	"cloudStorage/util"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

//分块上传初始化信息
type MultipartUploadInfo struct {
	FileHash   string //文件哈希
	FileSize   int    //文件大小
	UploadId   string //分块上传的唯一标识
	ChunkSize  int    //每个分块的大小
	ChunkCount int    //分块数量
}

//初始化分块上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	userid := r.Form.Get("userid")
	filehash := r.Form.Get("hash")
	fielsize, err := strconv.Atoi(r.Form.Get("size"))
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	fmt.Println(userid, filehash, fielsize)

	//获取redis连接
	redis := myredis.RedisPool().Get()
	defer redis.Close()

	//生成分块上传的初始化信息
	rand.Seed(time.Now().UnixNano())
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   fielsize,
		UploadId:   fmt.Sprintf("%s%x%d", userid, time.Now().UnixNano(), rand.Intn(1000)),
		ChunkSize:  10 * 1024, //5MB
		ChunkCount: int(math.Ceil(float64(fielsize) / (10 * 1024))),
	}

	//初始化信息写入redis
	redis.Do("HSET", "MP_"+upInfo.UploadId, "ChunkCount", upInfo.ChunkCount)
	redis.Do("HSET", "MP_"+upInfo.UploadId, "FileSize", upInfo.FileSize)
	redis.Do("HSET", "MP_"+upInfo.UploadId, "FileHash", upInfo.FileHash)

	//返回初始化信息
	util.SuccessResponse(w, upInfo)
}

//上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	uploadId := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	//获取redis连接
	rd := myredis.RedisPool().Get()
	defer rd.Close()

	//获得文件句柄，用于存储分块内容
	fpath := config.TEMP_CHUNK_PATH + uploadId + "/" + chunkIndex
	if err := os.MkdirAll(path.Dir(fpath), 0744); err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	fd, err := os.Create(fpath)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	//更新redis缓存状态
	rd.Do("HSET", "MP_"+uploadId, "chkidx_"+chunkIndex, 1)

	//返回处理结果到客户端
	util.SuccessResponse(w, nil)
}

//通知上传合并
func CompleteUploadHander(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	uploadId := r.Form.Get("uploadid")
	filename := r.Form.Get("filename")
	filehash := r.Form.Get("filehash")
	userid, err := strconv.Atoi(r.Form.Get("userid"))
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	//获取redis连接
	rd := myredis.RedisPool().Get()
	defer rd.Close()

	//通过uploadid查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rd.Do("HGETALL", "MP_"+uploadId))
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "ChunkCount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		util.ErrorResponse(w, "分块上传未完成")
		return
	}

	//合并分块

	//更新唯一文件表及用户文件表
	dao.OnFileUploadFinished(userid, filehash, filename, int64(filesize), "", false)

	//响应文件处理
	util.SuccessResponse(w, nil)
}
