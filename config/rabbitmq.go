package config

const (
	ASYNC_TRANSFER_ENABLE    = true                                 //是否开启文件异步转移
	RABBIT_URL               = "amqp://guest:guest@127.0.0.1:5672/" //rabbitmq的入口url
	TRANS_EXCHANGE_NAME      = "uploadserver.trans"                 //用于文件transfer的交换机
	TRANS_OSS_QUEUE_NAME     = "uploadserver.trans.oss"             //oss转移队列名
	TRANS_OSS_ERR_QUEUE_NAME = "uploadserver.trans.oss.err"         //oss转移失败后写入另一个队列的队列名
	TRANS_OSS_ROUTING_KEY    = "oss"                                //routing key
)
