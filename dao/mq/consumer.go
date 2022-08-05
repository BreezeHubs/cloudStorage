package mq

import "fmt"

//开始监听队列，获取消息
func StartConsumer(qName, cName string, callback func(msg []byte) bool) {
	fmt.Println(qName, cName, callback([]byte{}))
	//通过channel.Consumer获得消息信道
	msgs, err := channel.Consume(
		qName,
		cName,
		true,  //是否自动应答ack
		false, //指定是否唯一的消费者
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	//循环获取队列的消息

	done := make(chan bool)

	go func() {
		for d := range msgs {
			//调用callback函数处理消息
			if !callback(d.Body) {
				//如果callback函数返回false，则表示消息处理失败，则拒绝消息
				d.Nack(false, false)
			} else {
				//如果callback函数返回true，则表示消息处理成功，则确认消息
				d.Ack(false)
			}
		}
	}()

	<-done

	//关闭channel
	channel.Close()

}
